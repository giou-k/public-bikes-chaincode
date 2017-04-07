package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	"regexp"
)

var logger = shim.NewLogger("GiouChaincode")

//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
//CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4
const   AUTHORITY      	=  "regulator"
const   TOWNSHIP   		=  "township"
const   PRIVATE_ENTITY 	=  "private"
const   STATION  		=  "station"


//==============================================================================================================================
//	 Status types - Asset lifecycle is broken down into 4 statuses, this is part of the business logic to determine what can
//					be done to the bike at points in it's lifecycle
//==============================================================================================================================
const   STATE_TEMPLATE  			=  0
const   STATE_TOWNSHIP  			=  1
const   STATE_PRIVATE_OWNERSHIP 	=  2
const   STATE_STATION	 			=  3

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type  SimpleChaincode struct {
}

//==============================================================================================================================
//	Bike - Defines the structure for a bike object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Bike struct {
	Model           string `json:"model"`	
	Reg             string `json:"reg"`		//registration
	BTWSH           int    `json:"BTWSH"`	//"BikeTownship" : An integer value given by the Township that the bike belongs
	Owner           string `json:"owner"`
	Status          int    `json:"status"`	//values : 0(Regulator), 1(Township), 2(Private Entity), 3(Station)
	Colour          string `json:"colour"`
	crtID           string `json:"crtID"`	//"createID" : unique key
}


//==============================================================================================================================
//	crt Holder - Defines the structure that holds all the crtIDs for bikes that have been created.
//				Used as an index when querying all bikes.
//==============================================================================================================================

type crt_Holder struct {
	crts 	[]string `json:"crts"`
}

//==============================================================================================================================
//	User_and_eCert - Struct for storing the JSON of a user and their ecert
//==============================================================================================================================

type User_and_eCert struct {
	Identity string `json:"identity"`
	eCert string `json:"ecert"`
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var crtIDs crt_Holder

	bytes, err := json.Marshal(crtIDs)

    if err != nil { return nil, errors.New("Error creating crt_Holder record") }

	err = stub.PutState("crtIDs", bytes)

	for i:=0; i < len(args); i=i+2 {
		t.add_ecert(stub, args[i], args[i+1])
	}

	return nil, nil
}

//==============================================================================================================================
//	 General Functions
//==============================================================================================================================
//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
//				 for that user. Returns the ecert as retrived including html encoding.
//==============================================================================================================================
func (t *SimpleChaincode) get_ecert(stub shim.ChaincodeStubInterface, name string) ([]byte, error) {

	ecert, err := stub.GetState(name)

	if err != nil { return nil, errors.New("Couldn't retrieve ecert for user " + name) }

	return ecert, nil
}

//==============================================================================================================================
//	 add_ecert - Adds a new ecert and user pair to the table of ecerts
//==============================================================================================================================

func (t *SimpleChaincode) add_ecert(stub shim.ChaincodeStubInterface, name string, ecert string) ([]byte, error) {

	err := stub.PutState(name, []byte(ecert))

	if err == nil {
		return nil, errors.New("Error storing eCert for user " + name + " identity: " + ecert)
	}

	return nil, nil

}

//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================

func (t *SimpleChaincode) get_username(stub shim.ChaincodeStubInterface) (string, error) {

    username, err := stub.ReadCertAttribute("username");
	if err != nil { return "", errors.New("Couldn't get attribute 'username'. Error: " + err.Error()) }
	return string(username), nil
}

//==============================================================================================================================
//	 check_affiliation - Takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// 				  		certificates common name. The affiliation is stored as part of the common name.
//==============================================================================================================================

func (t *SimpleChaincode) check_affiliation(stub shim.ChaincodeStubInterface) (string, error) {
	
    affiliation, err := stub.ReadCertAttribute("role");
	if err != nil { return "", errors.New("Couldn't get attribute 'role'. Error: " + err.Error()) }
	return string(affiliation), nil

}

//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================

func (t *SimpleChaincode) get_caller_data(stub shim.ChaincodeStubInterface) (string, string, error){

	user, err := t.get_username(stub)

	affiliation, err := t.check_affiliation(stub);

    if err != nil { return "", "", err }

	return user, affiliation, nil
}

//==============================================================================================================================
//	 retrieve_crt - Gets the state of the data at crtID in the ledger then converts it from the stored
//					JSON into the Bike struct for use in the contract. Returns the Bike struct.
//					Returns empty v if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) retrieve_crt(stub shim.ChaincodeStubInterface, crtID string) (Bike, error) {

	var v Bike

	bytes, err := stub.GetState(crtID);

	if err != nil {	fmt.Printf("RETRIEVE_crt: Failed to invoke bike_code: %s", err); return v, errors.New("RETRIEVE_crt: Error retrieving bike with crtID = " + crtID) }

	err = json.Unmarshal(bytes, &v);

    if err != nil {	fmt.Printf("RETRIEVE_crt: Corrupt bike record "+string(bytes)+": %s", err); return v, errors.New("RETRIEVE_crt: Corrupt bike record"+string(bytes))	}

	return v, nil
}

//==============================================================================================================================
// save_changes - Writes to the ledger the Bike struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, v Bike) (bool, error) {

	bytes, err := json.Marshal(v)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error converting bike record: %s", err); return false, errors.New("Error converting bike record") }

	err = stub.PutState(v.crtID, bytes)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error storing bike record: %s", err); return false, errors.New("Error storing bike record") }

	return true, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	caller, caller_affiliation, err := t.get_caller_data(stub)

	if err != nil { return nil, errors.New("Error retrieving caller information")}


	if function == "create_bike" { return t.create_bike(stub, caller, caller_affiliation, args[0])
	} else { 							// If the function is not a create then there must be a bike so we need to retrieve the bike.

		argPos := 1

		v, err := t.retrieve_crt(stub, args[argPos])

        if err != nil { fmt.Printf("INVOKE: Error retrieving crt: %s", err); return nil, errors.New("Error retrieving crt") }


        if strings.Contains(function, "update") == false {

				if 		   function == "authority_to_township" { return t.authority_to_township(stub, v, caller, caller_affiliation, args[0], "township")
				} else if  function == "township_to_station"   { return t.township_to_station(stub, v, caller, caller_affiliation, args[0], "station")
				} else if  function == "station_to_private"		{ return t.station_to_private(stub, v, caller, caller_affiliation, args[0], "private")
				} else if  function == "private_to_station"  { return t.private_to_station(stub, v, caller, caller_affiliation, args[0], "station")
				}

		} else if function == "update_model"        { return t.update_model(stub, v, caller, caller_affiliation, args[0])
		} else if function == "update_reg" 			{ return t.update_registration(stub, v, caller, caller_affiliation, args[0])
		} else if function == "update_btwsh" 		{ return t.update_btwsh(stub, v, caller, caller_affiliation, args[0])
        } else if function == "update_colour" 		{ return t.update_colour(stub, v, caller, caller_affiliation, args[0])
		} 

														return nil, errors.New("Function of the name "+ function +" doesn't exist.")

	}
}
//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	caller, caller_affiliation, err := t.get_caller_data(stub)
	if err != nil { fmt.Printf("QUERY: Error retrieving caller details", err); return nil, errors.New("QUERY: Error retrieving caller details: "+err.Error()) }

    logger.Debug("function: ", function)
    logger.Debug("caller: ", caller)
    logger.Debug("affiliation: ", caller_affiliation)

	if function == "get_bike_details" {
		if len(args) != 1 { fmt.Printf("Incorrect number of arguments passed"); return nil, errors.New("QUERY: Incorrect number of arguments passed") }
		v, err := t.retrieve_crt(stub, args[0])
		if err != nil { fmt.Printf("QUERY: Error retrieving crt: %s", err); return nil, errors.New("QUERY: Error retrieving crt "+err.Error()) }
		return t.get_bike_details(stub, v, caller, caller_affiliation)
	} else if function == "check_unique_crt" {
		return t.check_unique_crt(stub, args[0], caller, caller_affiliation)
	} else if function == "get_bikes" {
		return t.get_bikes(stub, caller, caller_affiliation)
	} else if function == "get_ecert" {
		return t.get_ecert(stub, args[0])
	} else if function == "ping" {
        return []byte("Hello, world!"), nil
    }

	return nil, errors.New("Received unknown function invocation")

}

//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
//	 Create Bike - Creates the initial JSON for the bike and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_bike(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, crtID string) ([]byte, error) {
	var v Bike

	crt_ID         := "\"crtID\":\""+crtID+"\", "							// Variables to define the JSON
	btwsh          := "\"BTWSH\":0, "
	model          := "\"Model\":\"UNDEFINED\", "
	reg            := "\"Reg\":\"UNDEFINED\", "
	owner          := "\"Owner\":\""+caller+"\", "
	colour         := "\"Colour\":\"UNDEFINED\", "
	status         := "\"Status\":0 "

	bike_json := "{"+crt_ID+btwsh+model+reg+owner+colour+status+"}" 	// Concatenates the variables to create the total JSON object

	matched, err := regexp.Match("^[A-z][A-z][0-9]{7}", []byte(crtID))  // matched = true if the crtID passed fits format of two letters followed by seven digits

							if err != nil { fmt.Printf("create_bike: Invalid crtID: %s", err); return nil, errors.New("Invalid crtID") }

	if 				crt_ID  == "" 	 ||
					matched == false    {
																		fmt.Printf("create_bike: Invalid crtID provided");
																		return nil, errors.New("Invalid crtID provided")
	}

	err = json.Unmarshal([]byte(bike_json), &v)							// Convert the JSON defined above into a bike object for go

																		if err != nil { return nil, errors.New("Invalid JSON object") }

	record, err := stub.GetState(v.crtID) 								// If not an error then a record exists so cant create a new bike with this crtID as it must be unique

																		if record != nil { return nil, errors.New("Bike already exists") }

	if 	caller_affiliation != AUTHORITY {							// Only the regulator can create a new crt

		return nil, errors.New(fmt.Sprintf("Permission Denied. create_bike. %v === %v", caller_affiliation, AUTHORITY))

	}

	_, err  = t.save_changes(stub, v)

																		if err != nil { fmt.Printf("create_bike: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	bytes, err := stub.GetState("crtIDs")

																		if err != nil { return nil, errors.New("Unable to get crtIDs") }

	var crtIDs crt_Holder

	err = json.Unmarshal(bytes, &crtIDs)

																		if err != nil {	return nil, errors.New("Corrupt crt_Holder record") }

	crtIDs.crts = append(crtIDs.crts, crtID)


	bytes, err = json.Marshal(crtIDs)

															if err != nil { fmt.Print("Error creating crt_Holder record") }

	err = stub.PutState("crtIDs", bytes)

															if err != nil { return nil, errors.New("Unable to put the state") }

	return nil, nil

}

//=================================================================================================================================
//	 Transfer Functions
//=================================================================================================================================
//	 authority_to_township
//=================================================================================================================================
func (t *SimpleChaincode) authority_to_township(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if     	v.Status				== STATE_TEMPLATE	&&
			v.Owner					== caller			&&
			caller_affiliation		== AUTHORITY		&&
			recipient_affiliation	== TOWNSHIP			{		// If the roles and users are ok

					v.Owner  = recipient_name		// then make the owner the new owner
					v.Status = STATE_TOWNSHIP			// and mark it in the state of township

	} else {									// Otherwise if there is an error
															fmt.Printf("AUTHORITY_TO_TOWNSHIP: Permission Denied");
                                                            return nil, errors.New(fmt.Sprintf("Permission Denied. authority_to_township. %v %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_TEMPLATE, v.Owner, caller, caller_affiliation, AUTHORITY, recipient_affiliation, TOWNSHIP))


	}

	_, err := t.save_changes(stub, v)						// Write new state

															if err != nil {	fmt.Printf("AUTHORITY_TO_TOWNSHIP: Error saving changes: %s", err); return nil, errors.New("Error saving changes")	}

	return nil, nil									// We are Done

}

//=================================================================================================================================
//	 township_to_station
//=================================================================================================================================
func (t *SimpleChaincode) township_to_station(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if 		v.Model  == "UNDEFINED" ||
			v.Reg 	 == "UNDEFINED" ||
			v.Colour == "UNDEFINED" ||					
			v.BTWSH == 0				{			//If any part of the bike is undefined it has not been fully defined so cannot be sent
													fmt.Printf("TOWNSHIP_TO_PRIVATE: Bike not fully defined")
													return nil, errors.New(fmt.Sprintf("Bike not fully defined. %v", v))
	}

	if 		v.Status				== STATE_TOWNSHIP	&&
			v.Owner					== caller				&&
			caller_affiliation		== TOWNSHIP			&&
			recipient_affiliation	== STATION						{

					v.Owner = recipient_name
					v.Status = STATE_STATION

	} else {
        return nil, errors.New(fmt.Sprintf("Permission Denied. township_to_station. %v %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_TOWNSHIP, v.Owner, caller, caller_affiliation, TOWNSHIP, recipient_affiliation, STATION))
    }

	_, err := t.save_changes(stub, v)

	if err != nil { fmt.Printf("TOWNSHIP_TO_STATION: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 station_to_private
//=================================================================================================================================
func (t *SimpleChaincode) station_to_private(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if		v.Status				== STATE_STATION	&&
			v.Owner  				== caller					&&
			caller_affiliation		== STATION			&&
			recipient_affiliation	== PRIVATE_ENTITY			{

				v.Owner = recipient_name

	} else {
		return nil, errors.New(fmt.Sprintf("Permission Denied. station_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_STATION, v.Owner, caller, caller_affiliation, STATION, recipient_affiliation, PRIVATE_ENTITY))
	}

	_, err := t.save_changes(stub, v)
															if err != nil { fmt.Printf("STATION_TO_PRIVATE: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 private_to_station
//=================================================================================================================================
func (t *SimpleChaincode) private_to_station(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if 		v.Status				== STATE_PRIVATE_OWNERSHIP	&&
			v.Owner					== caller					&&
			caller_affiliation		== PRIVATE_ENTITY			&&
			recipient_affiliation	== STATION					{

					v.Owner = recipient_name

	} else {
        return nil, errors.New(fmt.Sprintf("Permission denied. private_to_station. %v === %v, %v === %v, %v === %v, %v === %v", v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, STATION))

	}

	_, err := t.save_changes(stub, v)
															if err != nil { fmt.Printf("PRIVATE_TO_STATION: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 Update Functions
//=================================================================================================================================
//	 update_btwsh
//=================================================================================================================================
func (t *SimpleChaincode) update_btwsh(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	new_btwsh, err := strconv.Atoi(string(new_value)) 		                // will return an error if the new btwsh contains non numerical chars

															if err != nil || len(string(new_value)) != 15 { return nil, errors.New("Invalid value passed for new BTWSH") }
	// Can't change the BTWSH after its initial assignment (BTWSH==0)
	if 		v.Status			== STATE_TOWNSHIP	&&
			v.Owner				== caller			&&
			caller_affiliation	== TOWNSHIP			&&
			v.BTWSH				== 0				{

					v.BTWSH = new_btwsh					// Update to the new value
	} else {

        return nil, errors.New(fmt.Sprintf("Permission denied. update_btwsh %v %v %v %v %v", v.Status, STATE_TOWNSHIP, v.Owner, caller, v.BTWSH))

	}

	_, err  = t.save_changes(stub, v)						// Save the changes in the blockchain

															if err != nil { fmt.Printf("UPDATE_BTWSH: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}


//=================================================================================================================================
//	 update_registration
//=================================================================================================================================
func (t *SimpleChaincode) update_registration(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, new_value string) ([]byte, error) {


	if		v.Owner	== caller	{

					v.Reg = new_value

	} else {
        return nil, errors.New(fmt.Sprint("Permission denied. update_registration"))
	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("UPDATE_REGISTRATION: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 update_colour
//=================================================================================================================================
func (t *SimpleChaincode) update_colour(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if 		v.Owner				== caller				&&
			caller_affiliation	== TOWNSHIP{

					v.Colour = new_value
	} else {

		return nil, errors.New(fmt.Sprint("Permission denied. update_colour %t %t" + v.Owner == caller, caller_affiliation == TOWNSHIP))
	}

	_, err := t.save_changes(stub, v)

		if err != nil { fmt.Printf("UPDATE_COLOUR: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 update_model
//=================================================================================================================================
func (t *SimpleChaincode) update_model(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if 		v.Status			== STATE_TOWNSHIP	&&
			v.Owner				== caller				&&
			caller_affiliation	== TOWNSHIP			{

					v.Model = new_value

	} else {
        return nil, errors.New(fmt.Sprint("Permission denied. update_model %t %t" + v.Owner == caller, caller_affiliation == TOWNSHIP))

	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("UPDATE_MODEL: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 Read Functions
//=================================================================================================================================
//	 get_bike_details
//=================================================================================================================================
func (t *SimpleChaincode) get_bike_details(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string) ([]byte, error) {

	bytes, err := json.Marshal(v)

																if err != nil { return nil, errors.New("GET_BIKE_DETAILS: Invalid bike object") }

	if 		v.Owner				== caller		||
			caller_affiliation	== AUTHORITY	{

					return bytes, nil
	} else {
																return nil, errors.New("Permission Denied. get_bike_details")
	}

}

//=================================================================================================================================
//	 get_bikes
//=================================================================================================================================

func (t *SimpleChaincode) get_bikes(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string) ([]byte, error) {
	bytes, err := stub.GetState("crtIDs")

																			if err != nil { return nil, errors.New("Unable to get crtIDs") }

	var crtIDs crt_Holder

	err = json.Unmarshal(bytes, &crtIDs)

																			if err != nil {	return nil, errors.New("Corrupt crt_Holder") }

	result := "["

	var temp []byte
	var v Bike

	for _, crt := range crtIDs.crts {

		v, err = t.retrieve_crt(stub, crt)

		if err != nil {return nil, errors.New("Failed to retrieve crt")}

		temp, err = t.get_bike_details(stub, v, caller, caller_affiliation)

		if err == nil {
			result += string(temp) + ","
		}
	}

	if len(result) == 1 {
		result = "[]"
	} else {
		result = result[:len(result)-1] + "]"
	}

	return []byte(result), nil
}

//=================================================================================================================================
//	 check_unique_crt
//=================================================================================================================================
func (t *SimpleChaincode) check_unique_crt(stub shim.ChaincodeStubInterface, crt string, caller string, caller_affiliation string) ([]byte, error) {
	_, err := t.retrieve_crt(stub, crt)
	if err == nil {
		return []byte("false"), errors.New("crt is not unique")
	} else {
		return []byte("true"), nil
	}
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {

	err := shim.Start(new(SimpleChaincode))

															if err != nil { fmt.Printf("Error starting Chaincode: %s", err) }
}

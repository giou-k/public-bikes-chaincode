package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	//allagh
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	// "crypto/x509"
	// "encoding/pem"
	// "net/url"
	"regexp"
)

var logger = shim.NewLogger("CLDChaincode")

//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
//CURRENT WORKAROUND USES ROLES CHANGE WHEN OWN USERS CAN BE CREATED SO THAT IT READ 1, 2, 3, 4, 5
const   AUTHORITY      	=  "regulator"
const   TOWNSHIP   		=  "township"
const   PRIVATE_ENTITY 	=  "private"
const   STATION  		=  "station"
const   SERVICE 		=  "service"


//==============================================================================================================================
//	 Status types - Asset lifecycle is broken down into 5 statuses, this is part of the business logic to determine what can
//					be done to the bike at points in it's lifecycle
//==============================================================================================================================
const   STATE_TEMPLATE  			=  0
const   STATE_TOWNSHIP  			=  1
const   STATE_PRIVATE_OWNERSHIP 	=  2
const   STATE_STATION	 			=  3
const   STATE_SERVICE		  		=  4

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type  SimpleChaincode struct {
}

//==============================================================================================================================
//	Bike - Defines the structure for a car object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Bike struct {
	Make            string `json:"make"`	//dn kserw se ti xrhsimopoieitai akoma!!
	Model           string `json:"model"`	//polhs,dromou,mtb
	Reg             string `json:"reg"`		//registration
	VIN             int    `json:"VIN"`
	Owner           string `json:"owner"`
	Serviced        bool   `json:"serviced"`//serviced or not
	Status          int    `json:"status"`	//pairnei times 0(Author),1(manufa),2(Priv-Dealers&LeaseComp&Leasee),4(Priv-Scrap)
	Colour          string `json:"colour"`
	V5cID           string `json:"v5cID"`
	LeaseContractID string `json:"leaseContractID"`// na to ksanatsekarw
}


//==============================================================================================================================
//	V5C Holder - Defines the structure that holds all the v5cIDs for bikes that have been created.
//				Used as an index when querying all bikes.
//==============================================================================================================================

type V5C_Holder struct {
	V5Cs 	[]string `json:"v5cs"`
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

	//Args
	//				0
	//			peer_address

	var v5cIDs V5C_Holder

	bytes, err := json.Marshal(v5cIDs)

    if err != nil { return nil, errors.New("Error creating V5C_Holder record") }

	err = stub.PutState("v5cIDs", bytes)

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

    // if err != nil { return "", "", err }

	// ecert, err := t.get_ecert(stub, user);

    // if err != nil { return "", "", err }

	affiliation, err := t.check_affiliation(stub);

    if err != nil { return "", "", err }

	return user, affiliation, nil
}

//==============================================================================================================================
//	 retrieve_v5c - Gets the state of the data at v5cID in the ledger then converts it from the stored
//					JSON into the Bike struct for use in the contract. Returns the Vehcile struct.
//					Returns empty v if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) retrieve_v5c(stub shim.ChaincodeStubInterface, v5cID string) (Bike, error) {

	var v Bike

	bytes, err := stub.GetState(v5cID);

	if err != nil {	fmt.Printf("RETRIEVE_V5C: Failed to invoke bike_code: %s", err); return v, errors.New("RETRIEVE_V5C: Error retrieving bike with v5cID = " + v5cID) }

	err = json.Unmarshal(bytes, &v);

    if err != nil {	fmt.Printf("RETRIEVE_V5C: Corrupt bike record "+string(bytes)+": %s", err); return v, errors.New("RETRIEVE_V5C: Corrupt bike record"+string(bytes))	}

	return v, nil
}

//==============================================================================================================================
// save_changes - Writes to the ledger the Bike struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (t *SimpleChaincode) save_changes(stub shim.ChaincodeStubInterface, v Bike) (bool, error) {

	bytes, err := json.Marshal(v)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error converting bike record: %s", err); return false, errors.New("Error converting bike record") }

	err = stub.PutState(v.V5cID, bytes)

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
	} else { 																				// If the function is not a create then there must be a car so we need to retrieve the car.

		argPos := 1

		if function == "service_bike" {																// If its a scrap bike then only two arguments are passed (no update value) all others have three arguments and the v5cID is expected in the last argument
			argPos = 0
		}

		v, err := t.retrieve_v5c(stub, args[argPos])

        if err != nil { fmt.Printf("INVOKE: Error retrieving v5c: %s", err); return nil, errors.New("Error retrieving v5c") }


        if strings.Contains(function, "update") == false && function != "service_bike"    { 									// If the function is not an update or a scrappage it must be a transfer so we need to get the ecert of the recipient.
				// ecert, err := t.get_ecert(stub, args[0]);

																		// if err != nil { return nil, err }

				if 		   function == "authority_to_township" { return t.authority_to_township(stub, v, caller, caller_affiliation, args[0], "township")
				} else if  function == "township_to_station"   { return t.township_to_station(stub, v, caller, caller_affiliation, args[0], "station")
				} else if  function == "station_to_private"		{ return t.station_to_private(stub, v, caller, caller_affiliation, args[0], "private")
				} else if  function == "private_to_station"  { return t.private_to_station(stub, v, caller, caller_affiliation, args[0], "station")
				} else if  function == "station_to_service"  { return t.station_to_service(stub, v, caller, caller_affiliation, args[0], "service")
				} //else if  function == "service_to_station" { return t.service_to_station(stub, v, caller, caller_affiliation, args[0], "station")
				//}

		} else if function == "update_make"  	    { return t.update_make(stub, v, caller, caller_affiliation, args[0])
		} else if function == "update_model"        { return t.update_model(stub, v, caller, caller_affiliation, args[0])
		} else if function == "update_reg" { return t.update_registration(stub, v, caller, caller_affiliation, args[0])
		} else if function == "update_vin" 			{ return t.update_vin(stub, v, caller, caller_affiliation, args[0])
        } else if function == "update_colour" 		{ return t.update_colour(stub, v, caller, caller_affiliation, args[0])
		} else if function == "service_bike" 		{ return t.service_bike(stub, v, caller, caller_affiliation) }

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
		v, err := t.retrieve_v5c(stub, args[0])
		if err != nil { fmt.Printf("QUERY: Error retrieving v5c: %s", err); return nil, errors.New("QUERY: Error retrieving v5c "+err.Error()) }
		return t.get_bike_details(stub, v, caller, caller_affiliation)
	} else if function == "check_unique_v5c" {
		return t.check_unique_v5c(stub, args[0], caller, caller_affiliation)
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
//	 Create Bike - Creates the initial JSON for the vehcile and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_bike(stub shim.ChaincodeStubInterface, caller string, caller_affiliation string, v5cID string) ([]byte, error) {
	var v Bike

	v5c_ID         := "\"v5cID\":\""+v5cID+"\", "							// Variables to define the JSON
	vin            := "\"VIN\":0, "
	make           := "\"Make\":\"UNDEFINED\", "
	model          := "\"Model\":\"UNDEFINED\", "
	reg            := "\"Reg\":\"UNDEFINED\", "
	owner          := "\"Owner\":\""+caller+"\", "
	colour         := "\"Colour\":\"UNDEFINED\", "
	leaseContract  := "\"LeaseContractID\":\"UNDEFINED\", "
	status         := "\"Status\":0, "
	serviced       := "\"Serviced\":false"

	bike_json := "{"+v5c_ID+vin+make+model+reg+owner+colour+leaseContract+status+serviced+"}" 	// Concatenates the variables to create the total JSON object

	matched, err := regexp.Match("^[A-z][A-z][0-9]{7}", []byte(v5cID))  				// matched = true if the v5cID passed fits format of two letters followed by seven digits

												if err != nil { fmt.Printf("create_bike: Invalid v5cID: %s", err); return nil, errors.New("Invalid v5cID") }

	if 				v5c_ID  == "" 	 ||
					matched == false    {
																		fmt.Printf("create_bike: Invalid v5cID provided");
																		return nil, errors.New("Invalid v5cID provided")
	}

	err = json.Unmarshal([]byte(bike_json), &v)							// Convert the JSON defined above into a bike object for go

																		if err != nil { return nil, errors.New("Invalid JSON object") }

	record, err := stub.GetState(v.V5cID) 								// If not an error then a record exists so cant create a new car with this V5cID as it must be unique

																		if record != nil { return nil, errors.New("Bike already exists") }

	if 	caller_affiliation != AUTHORITY {							// Only the regulator can create a new v5c

		return nil, errors.New(fmt.Sprintf("Permission Denied. create_bike. %v === %v", caller_affiliation, AUTHORITY))

	}

	_, err  = t.save_changes(stub, v)

																		if err != nil { fmt.Printf("create_bike: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	bytes, err := stub.GetState("v5cIDs")

																		if err != nil { return nil, errors.New("Unable to get v5cIDs") }

	var v5cIDs V5C_Holder

	err = json.Unmarshal(bytes, &v5cIDs)

																		if err != nil {	return nil, errors.New("Corrupt V5C_Holder record") }

	v5cIDs.V5Cs = append(v5cIDs.V5Cs, v5cID)


	bytes, err = json.Marshal(v5cIDs)

															if err != nil { fmt.Print("Error creating V5C_Holder record") }

	err = stub.PutState("v5cIDs", bytes)

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
			recipient_affiliation	== TOWNSHIP		&&
			v.Serviced				== false			{		// If the roles and users are ok

					v.Owner  = recipient_name		// then make the owner the new owner
					v.Status = STATE_TOWNSHIP			// and mark it in the state of manufacture

	} else {									// Otherwise if there is an error
															fmt.Printf("AUTHORITY_TO_TOWNSHIP: Permission Denied");
                                                            return nil, errors.New(fmt.Sprintf("Permission Denied. authority_to_township. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_TEMPLATE, v.Owner, caller, caller_affiliation, AUTHORITY, recipient_affiliation, TOWNSHIP, v.Serviced, false))


	}

	_, err := t.save_changes(stub, v)						// Write new state

															if err != nil {	fmt.Printf("AUTHORITY_TO_TOWNSHIP: Error saving changes: %s", err); return nil, errors.New("Error saving changes")	}

	return nil, nil									// We are Done

}

//=================================================================================================================================
//	 township_to_station
//=================================================================================================================================
func (t *SimpleChaincode) township_to_station(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

//	if 		v.Make 	 == "UNDEFINED" ||
//			v.Model  == "UNDEFINED" ||
//			v.Reg 	 == "UNDEFINED" ||
//			v.Colour == "UNDEFINED" ||					//If service == true, einai se service opote dn ginetai na metaferthei
//			v.VIN == 0				{					//If any part of the car is undefined it has not bene fully manufacturered so cannot be sent
//															fmt.Printf("TOWNSHIP_TO_PRIVATE: Car not fully defined")
//															return nil, errors.New(fmt.Sprintf("Car not fully defined. %v", v))
//	}

	if 		v.Status				== STATE_TOWNSHIP	&&
			v.Owner					== caller				&&
			caller_affiliation		== TOWNSHIP			&&
			recipient_affiliation	== STATION		&&
			v.Serviced     			== false							{

					v.Owner = recipient_name
					v.Status = STATE_STATION

	} else {
        return nil, errors.New(fmt.Sprintf("Permission Denied. township_to_station. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_STATION, v.Owner, caller, caller_affiliation, STATION, recipient_affiliation, STATION, v.Serviced, false))
    }

	_, err := t.save_changes(stub, v)

	if err != nil { fmt.Printf("TOWNSHIP_TO_STATION: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 private_to_private
//=================================================================================================================================
//func (t *SimpleChaincode) private_to_private(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {
//
//	if 		v.Status				== STATE_PRIVATE_OWNERSHIP	&&
//			v.Owner					== caller					&&
//			caller_affiliation		== PRIVATE_ENTITY			&&
//			recipient_affiliation	== PRIVATE_ENTITY			&&
//			v.Scrapped				== false					{
//
//					v.Owner = recipient_name
//
//	} else {
//        return nil, errors.New(fmt.Sprintf("Permission Denied. private_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, SERVICE, v.Scrapped, false))
//	}
//
//	_, err := t.save_changes(stub, v)
//
//															if err != nil { fmt.Printf("PRIVATE_TO_PRIVATE: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }
//
//	return nil, nil
//
//}

//=================================================================================================================================
//	 private_to_station
//=================================================================================================================================
func (t *SimpleChaincode) private_to_station(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if 		v.Status				== STATE_PRIVATE_OWNERSHIP	&&
			v.Owner					== caller					&&
			caller_affiliation		== PRIVATE_ENTITY			&&
			recipient_affiliation	== STATION					&&
            v.Serviced     			== false					{

					v.Owner = recipient_name

	} else {
        return nil, errors.New(fmt.Sprintf("Permission denied. private_to_station. %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v.Status, STATE_PRIVATE_OWNERSHIP, v.Owner, caller, caller_affiliation, PRIVATE_ENTITY, recipient_affiliation, STATION, v.Serviced, false))

	}

	_, err := t.save_changes(stub, v)
															if err != nil { fmt.Printf("PRIVATE_TO_STATION: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 station_to_private
//=================================================================================================================================
func (t *SimpleChaincode) station_to_private(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if		v.Status				== STATE_STATION	&&
			v.Owner  				== caller					&&
			caller_affiliation		== STATION			&&
			recipient_affiliation	== PRIVATE_ENTITY			&&
			v.Serviced				== false					{

				v.Owner = recipient_name

	} else {
		return nil, errors.New(fmt.Sprintf("Permission Denied. station_to_private. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_STATION, v.Owner, caller, caller_affiliation, STATION, recipient_affiliation, PRIVATE_ENTITY, v.Serviced, false))
	}

	_, err := t.save_changes(stub, v)
															if err != nil { fmt.Printf("STATION_TO_PRIVATE: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 station_to_service
//=================================================================================================================================
func (t *SimpleChaincode) station_to_service(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, recipient_name string, recipient_affiliation string) ([]byte, error) {

	if		v.Status				== STATE_STATION	&&
			v.Owner					== caller					&&
			caller_affiliation		== STATION			&&
			recipient_affiliation	== SERVICE			&&
			v.Serviced				== false					{

					v.Owner = recipient_name
					v.Status = STATE_SERVICE

	} else {
        return nil, errors.New(fmt.Sprintf("Permission Denied. station_to_service. %v %v === %v, %v === %v, %v === %v, %v === %v, %v === %v", v, v.Status, STATE_STATION, v.Owner, caller, caller_affiliation, STATION, recipient_affiliation, SERVICE, v.Serviced, false))
	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("STATION_TO_SERVICE: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 Update Functions
//=================================================================================================================================
//	 update_vin
//=================================================================================================================================
func (t *SimpleChaincode) update_vin(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	new_vin, err := strconv.Atoi(string(new_value)) 		                // will return an error if the new vin contains non numerical chars

															if err != nil || len(string(new_value)) != 15 { return nil, errors.New("Invalid value passed for new VIN") }

	if 		v.Status			== STATE_TOWNSHIP	&&
			v.Owner				== caller				&&
			caller_affiliation	== TOWNSHIP			&&
			v.VIN				== 0					&&			// Can't change the VIN after its initial assignment
			v.Serviced			== false				{

					v.VIN = new_vin					// Update to the new value
	} else {

        return nil, errors.New(fmt.Sprintf("Permission denied. update_vin %v %v %v %v %v", v.Status, STATE_TOWNSHIP, v.Owner, caller, v.VIN, v.Serviced))

	}

	_, err  = t.save_changes(stub, v)						// Save the changes in the blockchain

															if err != nil { fmt.Printf("UPDATE_VIN: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}


//=================================================================================================================================
//	 update_registration
//=================================================================================================================================
func (t *SimpleChaincode) update_registration(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, new_value string) ([]byte, error) {


	if		v.Owner				== caller			&&
			caller_affiliation	!= SERVICE	&&
			v.Serviced			== false			{

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
			caller_affiliation	== TOWNSHIP			&&/*((v.Owner				== caller			&&
			caller_affiliation	== TOWNSHIP)		||
			caller_affiliation	== AUTHORITY)			&&*/
			v.Serviced			== false				{

					v.Colour = new_value
	} else {

		return nil, errors.New(fmt.Sprint("Permission denied. update_colour %t %t %t" + v.Owner == caller, caller_affiliation == TOWNSHIP, v.Serviced))
	}

	_, err := t.save_changes(stub, v)

		if err != nil { fmt.Printf("UPDATE_COLOUR: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 update_make
//=================================================================================================================================
func (t *SimpleChaincode) update_make(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if 		v.Status			== STATE_TOWNSHIP	&&
			v.Owner				== caller				&&
			caller_affiliation	== TOWNSHIP			&&
			v.Serviced			== false				{

					v.Make = new_value
	} else {

        return nil, errors.New(fmt.Sprint("Permission denied. update_make %t %t %t" + v.Owner == caller, caller_affiliation == TOWNSHIP, v.Serviced))


	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("UPDATE_MAKE: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 update_model
//=================================================================================================================================
func (t *SimpleChaincode) update_model(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string, new_value string) ([]byte, error) {

	if 		v.Status			== STATE_TOWNSHIP	&&
			v.Owner				== caller				&&
			caller_affiliation	== TOWNSHIP			&&
			v.Serviced			== false				{

					v.Model = new_value

	} else {
        return nil, errors.New(fmt.Sprint("Permission denied. update_model %t %t %t" + v.Owner == caller, caller_affiliation == TOWNSHIP, v.Serviced))

	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("UPDATE_MODEL: Error saving changes: %s", err); return nil, errors.New("Error saving changes") }

	return nil, nil

}

//=================================================================================================================================
//	 service_bike
//=================================================================================================================================
func (t *SimpleChaincode) service_bike(stub shim.ChaincodeStubInterface, v Bike, caller string, caller_affiliation string) ([]byte, error) {

	if		v.Status			== STATE_SERVICE	&&
			v.Owner				== caller				&&
			caller_affiliation	== SERVICE		&&
			v.Serviced			== false				{

					v.Serviced = true

	} else {
		return nil, errors.New("Permission denied. service_bike")
	}

	_, err := t.save_changes(stub, v)

															if err != nil { fmt.Printf("service_bike: Error saving changes: %s", err); return nil, errors.New("service_bikerror saving changes") }

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
	bytes, err := stub.GetState("v5cIDs")

																			if err != nil { return nil, errors.New("Unable to get v5cIDs") }

	var v5cIDs V5C_Holder

	err = json.Unmarshal(bytes, &v5cIDs)

																			if err != nil {	return nil, errors.New("Corrupt V5C_Holder") }

	result := "["

	var temp []byte
	var v Bike

	for _, v5c := range v5cIDs.V5Cs {

		v, err = t.retrieve_v5c(stub, v5c)

		if err != nil {return nil, errors.New("Failed to retrieve V5C")}

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
//	 check_unique_v5c
//=================================================================================================================================
func (t *SimpleChaincode) check_unique_v5c(stub shim.ChaincodeStubInterface, v5c string, caller string, caller_affiliation string) ([]byte, error) {
	_, err := t.retrieve_v5c(stub, v5c)
	if err == nil {
		return []byte("false"), errors.New("V5C is not unique")
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

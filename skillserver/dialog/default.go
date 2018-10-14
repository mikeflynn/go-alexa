package dialog

// Type will indicate type of dialog interaction to be sent to the user.
type Type string

const (
	// Delegate will indicate that the Alexa service should continue the dialog ineraction.
	Delegate Type = "Dialog.Delegate"

	// ElicitSlot will indicate to the Alexa service that the specific slot should be elicited from the user.
	ElicitSlot Type = "Dialog.ElicitSlot"

	// ConfirmSlot indicates to the Alexa service that the slot value should be confirmed by the user.
	ConfirmSlot Type = "Dialog.ConfirmSlot"

	// ConfirmIntent indicates to the Alexa service that the complete intent should be confimed by the user.
	ConfirmIntent Type = "Dialog.ConfirmIntent"
)

const (
	// Started indicates that the dialog interaction has just begun.
	Started string = "STARTED"

	// InProgress indicates that the dialog interation is continuing.
	InProgress string = "IN_PROGRESS"

	// Completed indicates that the dialog interaction has finished.
	// The intent and slot confirmation status should be checked.
	Completed string = "COMPLETED"
)

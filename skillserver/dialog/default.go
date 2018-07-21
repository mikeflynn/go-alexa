package dialog

type Type string

const (
	DELEGATE       Type = "Dialog.Delegate"
	ELICIT_SLOT    Type = "Dialog.ElicitSlot"
	CONFIRM_SLOT   Type = "Dialog.ConfirmSlot"
	CONFIRM_INTENT Type = "Dialog.ConfirmIntent"
)

const (
	STARTED     string = "STARTED"
	IN_PROGRESS string = "IN_PROGRESS"
	COMPLETED   string = "COMPLETED"
)

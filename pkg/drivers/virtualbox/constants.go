package virtualbox

import "regexp"

const (
	vboxManageCmd = "VBoxManage"
)

var (
	regexVMNameUUID = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	regexVMInfoLine = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	regexColonLine  = regexp.MustCompile(`(.+):\s+(.*)`)
)

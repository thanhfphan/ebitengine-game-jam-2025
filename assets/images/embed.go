package images

import (
    _ "embed"
)

var (
    //go:embed backgrounds/main_bg.jpg
    MainBG_jpg []byte

    //go:embed backgrounds/tet.jpg
    PlayBG_jpg []byte

    //go:embed backgrounds/tablecard.png
    TableBG_png []byte

    //go:embed backgrounds/card_back.png
    CardBack_png []byte

    //go:embed ui/setting_icon.png
    SettingIcon_png []byte
)
package sounds

import (
    _ "embed"
)

var (
    //go:embed card-swipe.wav
    CardSwipe_wav []byte

    //go:embed poppop.wav
    PopPop_wav []byte

    //go:embed ding-effect.wav
    DingEffect_wav []byte

    //go:embed vietnam-bamboo-flute.ogg
    BambooFlute_ogg []byte
)
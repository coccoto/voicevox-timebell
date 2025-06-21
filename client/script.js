'use strict'

let cachedSpeakers = null

async function assembleHourList() {
    const elemHourList = document.getElementById('hourList')
    
    Array.from({length: 24}, (_, i) => i).forEach(i => {
        const elemInput = document.createElement('input')
        elemInput.type = 'checkbox'
        elemInput.value = i

        const elemLabel = document.createElement('label')
        elemLabel.textContent = i + "時"

        elemHourList.appendChild(elemInput)
        elemHourList.appendChild(elemLabel)
    })
}

async function assembleSpeakerList() {
    const speakers = await fetchSpeakers()
    const elemSpeakerList = document.getElementById('speakerList')
    
    for (const speaker of speakers) {
        const elemOption = document.createElement('option')
        elemOption.value = speaker.speaker_uuid
        elemOption.textContent = speaker.name

        elemSpeakerList.appendChild(elemOption)
    }
}

async function assembleStyleList() {
    const speakers = await fetchSpeakers()
    const elemSpeakerList = document.getElementById('speakerList')
    const elemStyleList = document.getElementById('styleList')
    // スタイルリストを初期化する
    elemStyleList.innerHTML = ''

    // 選択されたスピーカーのスタイルを取得する
    const selectedSpeakerUuid = elemSpeakerList.value
    const matchedSpeaker = speakers.find(speaker => speaker.speaker_uuid === selectedSpeakerUuid)    
    for (const styles of matchedSpeaker.styles) {
        const elemOption = document.createElement('option')
        elemOption.value = styles.id
        elemOption.textContent = styles.name

        elemStyleList.appendChild(elemOption)
    }
}

async function fetchSpeakers() {
    if (cachedSpeakers === null) {
        const result = await fetch('http://localhost:50021/speakers')
        cachedSpeakers = await result.json()
    }
    return cachedSpeakers
}

assembleHourList()
assembleSpeakerList()

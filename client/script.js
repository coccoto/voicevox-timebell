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
    if (selectedSpeakerUuid === '') {
        return
    }
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
        const result = await fetch(window.location.origin + ':8080/api/speakers')
        cachedSpeakers = await result.json()
    }
    return cachedSpeakers
}

async function onClickSaveButton() {
    const elemHourList = document.getElementById('hourList')
    const elemSpeakerList = document.getElementById('speakerList')
    const elemStyleList = document.getElementById('styleList')

    const selectedHourList = Array.from(elemHourList.querySelectorAll('input:checked')).map(input => input.value)
    const selectedStyleId = elemStyleList.value

    if (selectedStyleId === '') {
        alert("音声スタイルを選択してください。")
        return
    }

    const data = {
        hourList: selectedHourList,
        styleId: selectedStyleId,
    }
    try {
        const response = await fetch(window.location.origin + ':8080/api/config', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        alert("設定を保存しました。")

    } catch (error) {
        console.error(error)
    }
}

assembleHourList()
assembleSpeakerList()
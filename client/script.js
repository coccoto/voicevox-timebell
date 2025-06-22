'use strict'

let cachedSpeakers = null
let cachedConfig = null

async function fetchSpeakers() {
    if (cachedSpeakers === null) {
        const result = await fetch(window.location.origin + ':8080/api/speakers')
        cachedSpeakers = await result.json()
    }
    return cachedSpeakers
}

async function fetchConfig() {
    if (cachedConfig === null) {
        const result = await fetch(window.location.origin + ':8080/api/config-read')
        cachedConfig = await result.json()
    }
    return cachedConfig
}

async function assembleHourList() {
    const config = await fetchConfig()
    const elemHourList = document.getElementById('hourList')
    
    Array.from({length: 24}, (_, i) => i).forEach(i => {
        const elemInput = document.createElement('input')
        elemInput.type = 'checkbox'
        elemInput.value = i
        if (config && config.hourList && config.hourList.includes(i.toString())) {
            elemInput.checked = true
        }

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

async function checkedAllHourList(isChecked) {
    const elemHourList = document.getElementById('hourList')
    const checkboxList = elemHourList.querySelectorAll('input[type="checkbox"]')
    for (const checkbox of checkboxList) {
        checkbox.checked = isChecked
    }
}

async function buttonDisabled(isDisabled) {
    const elemRegisterButton = document.getElementById('registerButton')
    const elemTestPlayButton = document.getElementById('testPlayButton')

    elemRegisterButton.disabled = isDisabled
    elemTestPlayButton.disabled = isDisabled
}

// User interaction handlers

async function onChangeSpeakerList() {
    assembleStyleList() 
}

async function onClickRegisterButton() {
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
        buttonDisabled(true)

        const response = await fetch(window.location.origin + ':8080/api/config-register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })

    } catch (error) {
        console.error(error)
    } finally {
        buttonDisabled(false)
    }
}

async function onClickTestPlayButton() {
    try {
        buttonDisabled(true)
        const response = await fetch(window.location.origin + ':8080/api/alert')
        
        if (! response.ok) {
            throw new Error('リクエストに失敗しました')
        }
    } catch (error) {
        console.error(error)
        alert('エラーが発生しました: ' + error.message)
    } finally {
        buttonDisabled(false)
    }
}

function onClickCheckedButton() {
    checkedAllHourList(true)
}

function onClickUncheckedButton() {
    checkedAllHourList(false)
}

// init

window.addEventListener('DOMContentLoaded', function() {
    assembleHourList()
    assembleSpeakerList()
})

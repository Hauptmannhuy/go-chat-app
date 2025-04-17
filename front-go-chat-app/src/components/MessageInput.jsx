import { useState } from "react"
import { useUploadImage } from "../modules/useUploadImage"
import { FileImage } from "../types/image"

export function MessageInput({selectedChat, onSend}){
  const [inputValue, setInputValue] = useState("")
  const {inputImage, uploadFile} = useUploadImage()

  const sendMessage = () => {
    const payload = {input: null, image: null}
    
    if (inputImage != null) {
      payload.image = inputImage
    } 
     if (inputValue.length > 0) {
      payload.input = inputValue
    }
    console.log(inputImage)
    onSend(selectedChat, payload)
  }

  return (
    <div className="input-message">
      <button onClick={() => sendMessage()}>Send</button>
      <input type="file" onChange={(e) => {uploadFile(e.target.files[0])}}/>
      <input className="message-input" type="text" onChange={(e) => (setInputValue(e.target.value))} />
    </div>
  )
}
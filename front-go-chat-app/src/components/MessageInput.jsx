import { useState } from "react"

export function MessageInput({selectedChat, onSend}){
  const [inputValue, setInputValue] = useState("")

  return (
    <div className="input-message">
      <button onClick={() => onSend(selectedChat, inputValue)}>Send</button>
      <input className="message-input" type="text" onChange={(e) => (setInputValue(e.target.value))} />
    </div>
  )
}
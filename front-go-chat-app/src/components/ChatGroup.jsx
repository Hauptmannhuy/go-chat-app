import { useContext, useState } from "react"
import { GlobalContext } from "../contexts/GlobalContext"
import  { MessagesDisplay } from "../components/MessagesDisplay"

export default function ChatGroup({onSend}){
  const {selectedChat, messages, onlineStatus} = useContext(GlobalContext)
  const chatMessages = messages[selectedChat.name]
  const [inputValue, setInputValue] = useState("")
  console.log("chat changed to", selectedChat)
  console.log("chat messages", chatMessages)

  return (
    <>
      {(selectedChat.participation == false ) ? (
        <div>Start chatting by sending a message</div>
      )
    :
    (null)
    }
      < MessagesDisplay
      chatMessages={chatMessages}/>

      <input type="text" onChange={(e) => (setInputValue(e.target.value))} />
      <button onClick={() => onSend(selectedChat, inputValue)}>Send</button>
      </>
    )
  
}
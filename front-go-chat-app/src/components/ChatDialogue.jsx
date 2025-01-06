import { useContext, useState } from "react"
import { GlobalContext } from "../contexts/GlobalContext"
import { MessagesDisplay } from "./MessagesDisplay"
import { nameFormatter } from "../modules/nameFormatter"



function ChatDialogue({ onSend }){
  const {selectedChat, messages, onlineStatus} = useContext(GlobalContext)
  const chatMessages = messages[selectedChat.name]
  const [inputValue, setInputValue] = useState("")
  console.log("chat changed to", selectedChat)
  console.log("chat messages", chatMessages)

  return (
  <>
  <div>User status: {onlineStatus(nameFormatter(selectedChat.name))}</div>

   <MessagesDisplay
   chatMessages={chatMessages}/>
    
  <button onClick={() => onSend(selectedChat, inputValue)}>Send</button>
  <input type="text" onChange={(e) => (setInputValue(e.target.value))} />
     
      </>
    )
}

export default ChatDialogue
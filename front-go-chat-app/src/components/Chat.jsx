import { useContext, useState } from "react"
import { GlobalContext } from "../contexts/GlobalContext"
import { MessagesDisplay } from "./MessagesDisplay"



function Chat({ onSend }){
  const {selectedChat, messages} = useContext(GlobalContext)
  const chatMessages = messages[selectedChat.name]
  const [inputValue, setInputValue] = useState("")
  console.log("chat changed to", selectedChat)
  console.log("chat messages", chatMessages)

  const groupChatUnacquinted = () => ( selectedChat.chat_type == 'group' && !selectedChat.participation )
  return (
  <>
    {
      groupChatUnacquinted() ? (
        <button key={selectedChat.chat_id} onClick= {() => {subscribeHandler(selectedChat.chat_id) }}> Join {selectedChat.chat_id} group chat </button>
      ) : (
        null 
    )
  }

   <MessagesDisplay
   chatMessages={chatMessages}/>
    
  <button onClick={() => onSend(selectedChat, inputValue)}>Send</button>
  <input type="text" onChange={(e) => (setInputValue(e.target.value))} />
     
      </>
    )
}

export default Chat
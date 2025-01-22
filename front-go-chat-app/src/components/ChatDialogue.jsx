import { useContext, useState } from "react"
import { GlobalContext } from "../contexts/GlobalContext"
import { MessagesDisplay } from "./MessagesDisplay"
import { nameFormatter } from "../modules/nameFormatter"
import { MessageInput } from "./MessageInput"



function ChatDialogue({ onSend }){
  const {selectedChat, messages, onlineStatus} = useContext(GlobalContext)
  const chatMessages = messages[selectedChat.name]
  console.log("chat changed to", selectedChat)
  console.log("chat messages", chatMessages)

  return (
  <>
  <div className="message-display-container">
    <div>User status: {onlineStatus(nameFormatter(selectedChat.name))}</div>

   <MessagesDisplay
   chatMessages={chatMessages}/>
   <MessageInput
   onSend={onSend}
   selectedChat={selectedChat}/>
  </div>
      </>
    )
}

export default ChatDialogue
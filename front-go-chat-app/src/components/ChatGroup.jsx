import { useContext } from "react"
import { GlobalContext } from "../contexts/GlobalContext"
import  { MessagesDisplay } from "../components/MessagesDisplay"
import { MessageInput } from "./MessageInput"

export default function ChatGroup({onSend}){
  const {selectedChat, messages, onlineStatus} = useContext(GlobalContext)
  const chatMessages = messages[selectedChat.name]
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
      <div className="message-display-container">
      < MessagesDisplay
      chatMessages={chatMessages}/>

       <MessageInput
        onSend={onSend}
        selectedChat={selectedChat}/>
      </div>

      </>
    )
  
}
import { useContext } from "react"
import { ChatSnippet } from "./ChatSnippet"
import { GlobalContext } from "../contexts/GlobalContext"



function ChatList({onSelect, searchStatus}) {
  const {chats, searchResults} = useContext(GlobalContext)
  const currentChats = searchStatus ? searchResults : chats
  const chatKeys = Object.keys(currentChats)
  console.log(chats)
  return (
  <>
   <h2>Chat list</h2>
    <div className="chat-list">

    {chatKeys.map((key) => (
      < ChatSnippet 
      name={key}
      currentChats = {currentChats}
      onSelect = {onSelect}
      />
    )
  )}

  </div>
   
  
  
  </>
  )
}


export default ChatList
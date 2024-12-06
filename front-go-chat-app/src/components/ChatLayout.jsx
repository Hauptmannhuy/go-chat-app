import { useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";


import "../assets/chatBrowser.css"
import Chat from "./Chat";
import ChatList from "./ChatList";
import SignOutButton from "./SignOutButton";
import Search from "./SearchSection";
import { ChatContext } from "../contexts/ChatContext";
import SearchSection from "./SearchSection";
import { GlobalContext } from "../contexts/GlobalContext";
import {useEnvelope} from "../modules/makeEnvelope"


function ChatLayout(){ 

  const navigate = useNavigate()

  const {sendMessage, selectedChat, selectChat} = useContext(GlobalContext)
  const {makeEnvelope} = useEnvelope()

  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")
  
 

  function sendMessageHandler(chatObj, input){
    let message = null
    if (chatObj.type == 'private' && !chatObj.participation){ 
      let id = searchProfileResults[chatObj.name].id
       message = makeEnvelope("NEW_PRIVATE_CHAT", [id, input])
    } else {
      message = makeEnvelope("NEW_MESSAGE", [chatObj.name, input])
    }
    sendMessage(message)
  }
  
  const searchHandler = (input) => {
    if (input == "") return setEmptyInput()
    sendMessage(makeEnvelope("SEARCH_QUERY", [input, getUsernameCookie()]))
  }
 

  return (
  <>

    <SearchSection
    onSearch={searchHandler}
    />

    {selectedChat ? 
    (
      <Chat
      onSend={sendMessageHandler}/>
    )
    :
    (
      null
    )}

    

  </>
  )
}

export default ChatLayout
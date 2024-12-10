import { useContext, useState } from "react";
import { useNavigate } from "react-router-dom";


import "../assets/chatBrowser.css"
import Chat from "./Chat";
import ChatList from "./ChatList";
import SignOutButton from "./SignOutButton";
import SearchSection from "./SearchSection";
import { GlobalContext } from "../contexts/GlobalContext";
import {useEnvelope} from "../modules/makeEnvelope"
import { useAuth } from "../modules/useAuth";


function ChatLayout(){ 
  const navigate = useNavigate()

  const {sendMessage, selectedChat, selectChat, chats, searchProfileResults, searchResults} = useContext(GlobalContext)
  const {makeEnvelope} = useEnvelope()
  
  console.log(searchResults)
  console.log(chats)

  const [searchInputStatus, setSearchInputStatus] = useState(null)
  
  
  const {getUsername, userAuthenticated} = useAuth()

  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")

 
  
  const changeDialogueParticipation = () => {
    const modifiedChat = {...selectedChat, name: getUsername() + '_' + selectedChat.name, participation: true}
    selectChat(modifiedChat)
  }

  function sendMessageHandler(chatObj, input){
    let message = null
    if (chatObj.type == 'private' && !chatObj.participation){ 
      let id = searchProfileResults[chatObj.name].id
       message = makeEnvelope("NEW_PRIVATE_CHAT", [id, input])
       changeDialogueParticipation()
    } else {
      message = makeEnvelope("NEW_MESSAGE", [chatObj.name, input])
    }
    sendMessage(message)
  }
  
  const searchHandler = (input) => {
    if (input == "") return setSearchInputStatus(false)
    sendMessage(makeEnvelope("SEARCH_QUERY", [input, getUsername()]))
    setSearchInputStatus(true)
  }

  const selectChatHandler = (chat) => {
    selectChat(chat)
  }
 
  return (
  <>
  <SignOutButton />
    <SearchSection
    onSearch={searchHandler}
    />

    <ChatList
    onSelect={selectChatHandler}
    searchStatus={searchInputStatus}/>


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
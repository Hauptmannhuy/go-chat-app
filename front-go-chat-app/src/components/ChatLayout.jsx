import { useContext, useState } from "react";
import { useNavigate } from "react-router-dom";


import ChatDialogue from "./ChatDialogue";
import ChatList from "./ChatList";
import SignOutButton from "./SignOutButton";
import SearchSection from "./SearchSection";
import { GlobalContext } from "../contexts/GlobalContext";
import {useEnvelope} from "../modules/useEnvelope"
import { useAuth } from "../modules/useAuth";
import ChatGroup from "./ChatGroup";
import CreateGroupChat from "./GroupChatCreation";


function ChatLayout(){ 
  const navigate = useNavigate()

  const {writeToSocket, selectedChat, selectChat, chats, searchProfileResults, searchResults} = useContext(GlobalContext)
  const {makeEnvelope} = useEnvelope()
  
  const [searchInputStatus, setSearchInputStatus] = useState(null)
  const {getUsername, userAuthenticated} = useAuth()

  const changeDialogueParticipation = () => {
    const modifiedChat = {...selectedChat, name: getUsername() + '_' + selectedChat.name, participation: true}
    selectChat(modifiedChat)
  }



  /**
   * @param {{name: string, id: number, participation: boolean, type: string}} chatObj
   * @param {string} input 
   */

  function sendMessageHandler(chatObj, input){
    let message = null

    if (chatObj.type == 'private' && !chatObj.participation){ 
      let id = searchProfileResults[chatObj.name].id
      message = makeEnvelope("NEW_PRIVATE_CHAT", [id, input])
      changeDialogueParticipation()
    } else if (chatObj.type == 'group' && !chatObj.participation){
      message = makeEnvelope("JOIN_CHAT", [chatObj.id, chatObj.name, input])
    } else {
      message = makeEnvelope("NEW_MESSAGE", [chatObj.name, input])
    }
    writeToSocket(message)
  }
  function renderSelectedChat(){
    if (selectedChat.type == 'private'){
      return <ChatDialogue onSend={sendMessageHandler}/>
    } else {
      return <ChatGroup onSend={sendMessageHandler}/>
    }
  } 
  
  const searchHandler = (input) => {
    if (input == "") return setSearchInputStatus(false)
    writeToSocket(makeEnvelope("SEARCH_QUERY", [input, getUsername()]))
    setSearchInputStatus(true)
  }

  const selectChatHandler = (chat) => {
    selectChat(chat)
  }
 
  return (
  <div className="container">

  <div className="left-part">
    <SignOutButton />
      <SearchSection
      onSearch={searchHandler}
      />
      <CreateGroupChat/>
      <ChatList
      onSelect={selectChatHandler}
      searchStatus={searchInputStatus}/>
  </div>

    
    <div className="right-part">
      {selectedChat ?
      (
        renderSelectedChat()
      )
      :
      (
        null
      )}
    </div>

    

    </div>
  )
}

export default ChatLayout
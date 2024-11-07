import { useState } from "react";
import { useNavigate } from "react-router-dom";


import Chat from "./Chat";
import ChatList from "./ChatList";
import SignOutButton from "./SignOutButton";
import Search from "./Search";

import { useChatAction } from "../modules/useChatActions";
import { useWebsocket } from "../modules/useWebsocket";
import { useChatBuild } from "../modules/useChatBuild";
import { useMessageBuild } from "../modules/useMessageBuild";
import { useSearchQuery } from "../modules/useSearchQuery";

function ChatBrowser(){ 
  const navigate = useNavigate()
  const [chatSelected, setSelectedChat] = useState(null)
  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")
  
  const {messages, addMessageStorage, handleMessageLoad, addMessage} = useMessageBuild()
  const {searchResults, handleSearchQuery, setEmptyInput} = useSearchQuery()
  const {chats, handleInitChatLoad, addChat} = useChatBuild(addMessageStorage)

  const {sendMessage} = useWebsocket("/socket/chat", (ev) => {
    processSocketMessage(ev)
  })

  const {sendEnvelope, processSocketMessage} = useChatAction(sendMessage, {
    handleInitChatLoad, addChat, handleMessageLoad, addMessage, handleSearchQuery
  })


  

  function selectChatHandler(chatname) {
    console.log(chatname)
    if (!messages[chatname]) {
      addMessageStorage(chatname)
    }
    const selectedChat = chats[chatname] || searchResults[chatname]
    console.log(selectedChat)
    setSelectedChat(selectedChat)
  }
  
  console.log(messages)
    

  const createGroupChat = (name) => {
    setInvokeStatus(false)
    setCreationInput(false)
    
    sendEnvelope("NEW_CHAT", [name])
  };


  function MessageSendHandler(chatObj, input){
    console.log(chatObj)
    if (chatObj.type == 'private' && !chatObj.participation){ 
      console.log(searchProfileResults)
      let id = searchProfileResults[chatObj.name].profile.id
      
      return sendEnvelope("NEW_PRIVATE_CHAT", [id, input])
    }
    sendEnvelope("NEW_MESSAGE", [chatObj.name, input])
  }
  
 
  const search = (input) => {
    if (input == "") return setEmptyInput()
    sendEnvelope("SEARCH_QUERY", [input, getUsernameCookie()])
  }



  const userAuthenticated = () => {
    if (document.cookie != '') return false;
    return true;
  }

  const getUsernameCookie = () =>  document.cookie.split('=')[1]

  

  const joinChat = (name,type = "JOIN_CHAT") => {
    appendNewChat(name)
    sendEnvelope(type, [name, getUsernameCookie() ])
  }
  
  const appendNewChat = (name, type) => {
    const previuosChats = [...chats]
    const newChats = previuosChats.filter(el => (el.name != name))
    const changedChat = {name: name, participation: true}
    newChats.push(changedChat)
    setChats(newChats)
  }

  console.log(chatSelected)
  return (
    <>

    {userAuthenticated ? (< SignOutButton/>) : null}

    <div className="search-section">

    <p>
    <Search
    searchHandler={search}/>
    </p>

   { searchResults ?  
    <ChatList 
     chats={searchResults}
     handleSelect={selectChatHandler}
     /> : 
     <ChatList
      chats={chats}
      handleSelect={selectChatHandler}
      />
   } 
    </div>

    <div className="chat-bar">

    

    <button onClick={() => setInvokeStatus(true)}>Create chat</button>
    
    {creationChatInvoked ? (
      <div>
        <input type="text" onChange={(e) => setCreationInput(e.target.value)}/>
        <button onClick={() => createGroupChat(creationInput)}>Create</button>
      </div>
    ) : (
      null
    )}
    </div>
    
    <div className="chat-display">
      {chatSelected ? (
      <Chat
      chat = {chatSelected}
      msgHandler = {MessageSendHandler}
      messages = {messages[chatSelected.name]}
      subscribeHandler={joinChat}
      userID = {getUsernameCookie()}/>
      ) : (
      <h2>Chat display</h2>
          )}
    </div>
     
  </>
  )
}

export default ChatBrowser
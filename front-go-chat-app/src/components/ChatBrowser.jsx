import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";


import "../assets/chatBrowser.css"
import Chat from "./Chat";
import ChatList from "./ChatList";
import SignOutButton from "./SignOutButton";
import Search from "./Search";

import { useChatAction } from "../modules/useChatActions";
import { useWebsocket } from "../modules/useWebsocket";
import { useChatBuild } from "../modules/useChatBuild";
import { useMessageBuild } from "../modules/useMessageBuild";
import { useSearchQuery } from "../modules/useSearchQuery";
import { useIndexedDB } from "../modules/useIndexedDB";
function ChatBrowser(){ 

  const userAuthenticated = () => {
    if (document.cookie != '') return false;
    return true;
  }
  const getUsernameCookie = () =>  document.cookie.split('=')[1]

  const {messages, addMessageStorage, handleMessageLoad, addMessage} = useMessageBuild()
  const {searchResults,searchProfileResults, handleSearchQuery, setEmptyInput} = useSearchQuery()
  const {chats, handleInitChatLoad, addChat} = useChatBuild(addMessageStorage, getUsernameCookie())

  

  const {sendMessage, socket} = useWebsocket("/socket/chat", (ev) => {
    processSocketMessage(ev)
  })
  
  const {cacheStatus, cacheMessages} = useIndexedDB()

  const {sendEnvelope, processSocketMessage} = useChatAction(sendMessage, {
    handleInitChatLoad, addChat, cacheMessages, addMessage, handleSearchQuery
  })

   
 
  console.log(cacheStatus)
  const navigate = useNavigate()
  const [chatSelected, setSelectedChat] = useState(null)
  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")
  
 


  function selectChatHandler(chatname) {
    if (!messages[chatname]) {
      addMessageStorage(chatname)
    }
    const selectedChat = chats[chatname] || searchResults[chatname]
    setSelectedChat(selectedChat)
  }
  
    

  const createGroupChat = (name) => {
    setInvokeStatus(false)
    setCreationInput(false)
    
    sendEnvelope("NEW_CHAT", [name])
  };


  function MessageSendHandler(chatObj, input){
    if (chatObj.type == 'private' && !chatObj.participation){ 
      let id = searchProfileResults[chatObj.name].id
      return sendEnvelope("NEW_PRIVATE_CHAT", [id, input])
    }
    sendEnvelope("NEW_MESSAGE", [chatObj.name, input])
  }
  
 
  const search = (input) => {
    if (input == "") return setEmptyInput()
    sendEnvelope("SEARCH_QUERY", [input, getUsernameCookie()])
  }


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

  useEffect(() => {
    console.log("123")
    console.log(socket.current)
    if (socket.current != null) {
      console.log("FETCH MESSAGES")
      const message = {type: "LOAD_MESSAGES"}
      sendMessage(message)
    } 
  }, [socket.current])

  return (
    <>

    {userAuthenticated ? (< SignOutButton/>) : null}

    <button onClick={() => setInvokeStatus(true)}>Create chat</button>
    
    {creationChatInvoked ? (
      <div>
        <input type="text" onChange={(e) => setCreationInput(e.target.value)}/>
        <button onClick={() => createGroupChat(creationInput)}>Create</button>
      </div>
    ) : (
      null
    )}



      <div className="container">

        <div className="left-part">


    <div className="search-bar">
    <Search
    searchHandler={search}/>
    </div>


    <div className="chat-tab section">

   { searchResults ?  
    <ChatList 
     chats={searchResults}
     handleSelect={selectChatHandler}
     currentUsername={getUsernameCookie()}
     messages={messages}

     /> : 
     <ChatList
      chats={chats}
      handleSelect={selectChatHandler}
      currentUsername={getUsernameCookie()}
      messages={messages}

      />
   } 
    </div>
    </div>


    
   <div className="right-part">

  
    
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
      </div>

    </div>

     
  </>
  )
}

export default ChatBrowser
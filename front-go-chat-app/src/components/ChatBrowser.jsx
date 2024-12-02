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
  
  const { cacheMessages, cacheChats, getMessages, getChats, saveMessage} = useIndexedDB(
    () => {
      fetchCache()
    },
    () => {
      display()
    }
  )

  const {sendEnvelope, processSocketMessage, checkFetchStatus} = useChatAction(sendMessage, {
    cacheChats, addChat, cacheMessages, addMessage, handleSearchQuery, saveMessage
  })

   
 
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
  
  // const appendNewChat = (name, type) => {
  //   const previuosChats = [...chats]
  //   const newChats = previuosChats.filter(el => (el.name != name))
  //   const changedChat = {name: name, participation: true}
  //   newChats.push(changedChat)
  //   setChats(newChats)
  // }


  
  function fetchCache() {
    console.log("asking cache")
    sendMessage({type: "LOAD_SUBS"})
    sendMessage({type: "LOAD_MESSAGES"})
    checkFetchStatus()
    .then(() => {
      display()
    })
    .catch((reason) => {
      console.error("Error during cache fetch:", reason)
    })
  }

  function display() {
    const {privateChatReq, groupChatReq} = getChats()
    const messageReq = getMessages()

    privateChatReq.addEventListener("success", () => {
      
      const result = privateChatReq.result
      handleInitChatLoad(result, "private", true)
    })

    groupChatReq.addEventListener("success", () => {
      const result = groupChatReq.result
      
      handleInitChatLoad(result, "group", true)
    })

    
    messageReq.addEventListener("success", () => {
      handleMessageLoad(messageReq.result)
    })

  }
 
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
import { useEffect, useState } from "react";

import { useWebsocket } from "../modules/useWebsocket";
import { useChatAction } from "../modules/useChatActions";
import Chat from "./Chat";
import ChatList from "./ChatList";
import { useNavigate } from "react-router-dom";
import SignOutButton from "./SignOutButton";
import { useRef } from "react";
import Search from "./Search";

function ChatBrowser(){ 
  const navigate = useNavigate()
  const [chatSelected, setSelectedChat] = useState(null)
  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")
  
  const [profiles, setProfiles] = useState({})
  const [chats, setChats] = useState({});
  const [messages, setMessages] = useState({})
  const [searchResults, setSearchResults] = useState(null)
  const [searchProfileResults, setsearchProfileResults] = useState(null)
  

  const {sendMessage} = useWebsocket("/socket/chat", (ev) => {
    processSocketMessage(ev)
  })
  const {sendEnvelope, processSocketMessage} = useChatAction(sendMessage, {
    setChats, setMessages, setProfiles, setsearchProfileResults,setSearchResults
  })
  

  function selectChatHandler(chatname) {
    console.log(chatname)
    if (!messages[chatname]) {
      addMessagesObjectHandler(chatname)
    }
    const selectedChat = chats[chatname] || searchResults[chatname]
    console.log(selectedChat)
    setSelectedChat(selectedChat)
  }

    

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
    if (input == "") return setSearchResults(null)
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
import { useEffect, useState } from "react";

import Chat from "./Chat";
import ChatList from "./ChatList";
import { useNavigate } from "react-router-dom";
import SignOutButton from "./SignOutButton";
import { useRef } from "react";

function ChatBrowser(){ 
   

  const socketConnection = useRef(null)

  const navigate = useNavigate()
  
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState({})


  const [chatSelected, setSelectedChat] = useState(null)
  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")

  const sendEnvelope = (type, data) => {
      const actions = {
        "NEW_MESSAGE": {
          type: "NEW_MESSAGE",
          user_id: `${data[0]}`,
          chat_id: `${data[1]}`,
          body: `${data[2]}`,
        },
        "NEW_CHAT": {
          type: "NEW_CHAT",
          id: `${data[0]}`
        },
        "JOIN_CHAT": {
          type: "JOIN_CHAT",
          chat_id: `${data[0]}`,
          user_id: `${data[1]}`
        }
      }
    
      const json = actions[type];
      if (json) {
        socketConnection.current.send(JSON.stringify(json));
      } else {
        console.error(`Unknown action type: ${type}`);
      }
    };

  const createChat = (name,type = "NEW_CHAT") => {
    setInvokeStatus(false)
    setCreationInput(false)
    
    sendEnvelope(type, [name])
  };

  const joinChat = (name,type = "JOIN_CHAT") => {
   const previuosChats = [...chats]
   const newChats = previuosChats.filter(el => (el.name != name))
   const changedChat = {name: name, participation: true}
   newChats.push(changedChat)

   sendEnvelope(type, [name, getUsernameCookie() ])
   setChats(newChats)
  }

  const addChatsAndMessages = (names, participation = false) => {
    if (!names) return
    typeof names != 'object' ? names = [names] : null
    names.forEach(name => {
      addChatHandler(name, participation)
      addMessagesObjectHandler(name)
    })
  }

  const addChatHandler = (name, participation) => {
    setChats((chats) => {
      const newChat = {name: `${name}`, participation: participation}
      return [...chats, newChat]
    })
  }
  const addMessagesObjectHandler = (name) => {
    setMessages((messages) => {
      const newMessages = {...messages}
      newMessages[name] = []
      return newMessages
    })
  }

  function saveLocalMessage(message) {
    setMessages((messages) => {
      const newMessages = {...messages}
      console.log(newMessages)
      
    
      newMessages[message.chat_id].push(message)
      return newMessages
    })
  }

  function handleMessageLoad(data){
   setMessages((messages) => ({...messages, ...data}))
  }


  const sendMessage = (chatID, userID, msg, type = "NEW_MESSAGE") => {
    sendEnvelope(type, [userID, chatID, msg])
  }
  
 
  function processSocketMessage(ev) {
    const response = JSON.parse(ev.data);
    console.log(response.Data)
    switch (response.Type) {
      case "NEW_CHAT":
        addChatsAndMessages(response.Data.id);
        break;
      case "NEW_MESSAGE":
        saveLocalMessage(response.Data)
        break;
      case "LOAD_SUBS":
        return addChatsAndMessages(response.Data, true);
      case "LOAD_MESSAGES":
        handleMessageLoad(response.Data)
        console.log(response)
        break;
      default:
        console.log(ev.data);
        break;
    }
  }; 

  const userAuthenticated = () => {
    if (document.cookie != '') return false;
    return true;
  }

  const getUsernameCookie = () => {
    return document.cookie.split('=')[1]
   }
 
  useEffect(()=>{
    const websocket = new WebSocket("/socket/chat")

    websocket.addEventListener("open", () => {
      socketConnection.current = websocket
    })

    websocket.addEventListener("message", processSocketMessage)

    
    websocket.addEventListener("close", (ev) => {
     
    })
    return () => websocket.close()
  }, [])

  

  console.log("messages", messages)

  
  return (
    <>

    {userAuthenticated ? (< SignOutButton/>) : null}

    <div className="chat-bar">

     <ChatList 
     chats={chats}
     handleSelect={setSelectedChat}
     handleJoin={joinChat}
     />

    <button onClick={() => setInvokeStatus(true)}>Create chat</button>
    
    {creationChatInvoked ? (
      <div>
        <input type="text" onChange={(e) => setCreationInput(e.target.value)}/>
        <button onClick={() => createChat(creationInput)}>Create</button>
      </div>
    ) : (
      null
    )}
    </div>
    
    <div className="chat-display">
      {chatSelected ? (
      <Chat
      chatName = {chatSelected}
      msgHandler = {sendMessage}
      messages = {messages[chatSelected]}
      userID = {getUsernameCookie()}/>
      ) : (
      <h2>Chat display</h2>
          )}
    </div>
     
  </>
  )
}

export default ChatBrowser
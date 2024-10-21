import { useEffect, useState } from "react";

import Chat from "./Chat";
import ChatList from "./ChatList";
import { useNavigate } from "react-router-dom";
import SignOutButton from "./SignOutButton";

function ChatBrowser(){ 
  const navigate = useNavigate()
  
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState({})
  const [websocket, setWebsocket] = useState(null)

  const [chatSelected, setSelectedChat] = useState(null)
  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")

  const sendEnvelope = (type, data) => {
      const actions = {
        "NEW_MESSAGE": {
          type: "NEW_MESSAGE",
          chatid: `${data[0]}`,
          body: `${data[1]}`,
        },
        "NEW_CHAT": {
          type: "NEW_CHAT",
          id: `${data[0]}`
        },
        "JOIN_CHAT": {
          type: "JOIN_CHAT",
          chatID: `${data[0]}`,
          userID: `${data[1]}`
        }
      }
    
      const json = actions[type];
      if (json && websocket) {
        websocket.send(JSON.stringify(json));
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

  const appendChats = (names, participation = false) => {
    typeof names != 'object' ? names = [names] : null
    const newMessages = {...messages}
    const newChats = []
    names.forEach(name => {
      
     newChats.push({ name: `${name}`, participation: participation });
     newMessages[name] = []
    })
    setChats([...chats, ...newChats]);
    setMessages(newMessages)
  }
  console.log(messages)
  console.log(chats)

  const sendMessage = (name, msg, type = "NEW_MESSAGE") => {
    sendEnvelope(type, [name,msg])
  }
  
  if (websocket != null && websocket != "error") {
    websocket.onmessage = (ev) => {
      const response = JSON.parse(ev.data)
      console.log(response)
      switch (response.Type) {
        case "NEW_CHAT":
          appendChats(response.Data.id)
          break;
          case "GET_ALL_MESSAGES":
            console.log('All messages',ev.data)
            break
          case "LOAD_SUBS":
           return appendChats(response.Data, true)
          default:
            console.log(ev.data)
          break;
        }
  }
}






  

  useEffect(() => {

        const ws = new WebSocket(`/socket/chat`)
      
        ws.onopen = function () {
          console.log("WebSocket connection established.");
          console.log(ws.readyState, 'readyState')
          setWebsocket(ws)
        };
        ws.onerror = function (error) {
          console.error("WebSocket error:", error);
          setWebsocket("error")
      };
      ws.onclose = function (event) {
        if (event.code === 1006) {
            console.error("Connection closed abnormally, possibly due to redirection.");
           
        } else {
            console.log("WebSocket connection closed:", event);
        }
    }

    return () => {
      ws.close()
    }
  }, [])




  // if (websocket == "error"){
  //   navigate(`/sign_up`)
  // }

  const userAuthenticated = () => {
    if (document.cookie != '') return false;
    return true;
  }

  const getUsernameCookie = () => {
    return document.cookie.split('=')[1]
   }
 
  
  console.log(document.cookie)

  
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
      ws = {websocket}
      chatName = {chatSelected}
      msgHandler = {sendMessage}
      messages = {messages[chatSelected]}/>
      ) : (
      <h2>Chat display</h2>
          )}
    </div>
     
  </>
  )
}

export default ChatBrowser
import { useEffect, useState } from "react";
import Chat from "./Chat";
import ChatList from "./ChatList";
function ChatBrowser(){
  
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState({})
  const [websocket, setWebsocket] = useState(null)
  

  const [chatSelected, setSelectedChat] = useState(null)
  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")

  const sendEnvelope = (type, data) => {
    console.log(type)
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
          chatID: `${data[0]}`
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

   sendEnvelope(type, [name])
   setChats(newChats)
  }

  const appendChats = (name) => {
    const newChat = { name: `${name}`, participation: false };
    setChats([...chats, newChat]);
    const newMessages = {...messages}
    newMessages[name] = []
    setMessages(newMessages)
  }


  const sendMessage = (name, msg, type = "NEW_MESSAGE") => {
    sendEnvelope(type, [name,msg])
  }
  
  if (websocket) {
    websocket.onmessage = (ev) => {
      console.log(ev.data)
      const response = JSON.parse(ev.data)
      console.log(response)
      if (response.Type == 'NEW_CHAT'){
        appendChats(response.Data.ID)
      }
      else {
        console.log(ev.data)
      }
      
    }
  }
  console.log(chats,messages)

  

  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8090/chat`);
    ws.onopen = function () {
      console.log("WebSocket connection established.");
      console.log(ws.readyState, 'readyState')
    };
    setWebsocket(ws)

    return () => {
      ws.close()
    }

  }, [])
  
  

  
  return (
    <>
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
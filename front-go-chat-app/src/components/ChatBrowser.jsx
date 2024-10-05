import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Chat from "./Chat";
function ChatBrowser(){
  
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState({})
  const [websocket, setWebsocket] = useState(null)
  

  const [chatSelected, setSelectedChat] = useState(null)
  const [creationChatInvoked, setInvokeStatus] = useState(false)
  const [creationInput, setCreationInput] = useState("")


  const createChat = (name) => {
    setInvokeStatus(false)
    setCreationInput(false)
    
    sendChatData(name)
  };

  const joinChat = (name) => {
   const previuosChats = [...chats]
   const newChats = previuosChats.filter(el => (el.name != name))
   const changedChat = {name: name, participation: true}
   newChats.push(changedChat)

   const json_message = {
    "type": "JOIN_CHAT",
    "chatID": `${name}`
   }

   websocket.send(JSON.stringify(json_message))
   setChats(newChats)
  }

  const appendChats = (name) => {
    const newChat = { name: `${name}`, participation: false };
    setChats([...chats, newChat]);
    const newMessages = {...messages}
    newMessages[name] = []
    setMessages(newMessages)
  }
  // console.log(messages)

  const sendChatData = (name) => {
   const json_message = {
      "type": "NEW_CHAT",
      "id": `${name}`,
      }
    websocket.send(JSON.stringify(json_message))
  }

  const appendNewMessages = (name, msg) => {
    const newArr = {...messages}
    newArr[name].push(msg)
    setMessages(newArr)
  }
  
  if (websocket) {
    websocket.onmessage = (ev) => {
      console.log(ev.data)
      const response = JSON.parse(ev.data)
      console.log(response)
      if (response.Type == 'NEW_CHAT'){
        appendChats(response.Data.ID)
        console.log(response.ID)
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
      <h2>Chat list</h2>
    <div className="chat-list">
    {chats.map((el) => (
       el.participation ? ( 
        <button  

          key={el.name} 
          onClick= {() => {setSelectedChat(el.name) }}
          name={el.name}
          >{el.name} 

        </button>)
         : 

        (
        <button
        key={el.name}
        onClick = {() => {joinChat(el.name)}}>
          Join chat {el.name}
        </button>
        )

       
      
      ))}
    </div>
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
      msgHandler = {appendNewMessages}
      messages = {messages[chatSelected]}/>
      ) : (
      <h2>Chat display</h2>
          )}
    </div>
     
  </>
  )
}

export default ChatBrowser
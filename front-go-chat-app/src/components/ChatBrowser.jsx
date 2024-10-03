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
    const newChat = { name: `${name}` };
    setChats([...chats, newChat]);
    const newMessages = {...messages}
    newMessages[name] = []
    setMessages(newMessages)
    sendChatData(name)
  };
  console.log(messages)

  const sendChatData = (name) => {
   const json_message = {
      "type": "NEW_CHAT",
      "data": {
        "name": `${name}`
      }
    }
    websocket.send(JSON.stringify(json_message))
  }

  const appendNewMessages = (name, msg) => {
    const newArr = {...messages}
    newArr[name].push(msg)
    setMessages(newArr)
  }
  



  

  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8090/chat`);
    ws.onopen = function () {
      console.log("WebSocket connection established.");
      ws.send("Hello from client!");
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
      <button  
      key={el.name} 
      onClick= {() => {setSelectedChat(el.name) }}
      name={el.name}
      >{el.name} </button>
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
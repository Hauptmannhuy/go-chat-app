import { useEffect, useState } from "react";


function ChatRoom(){
  const [messages, setMessages] = useState([])
  const [inputValue, setInputValue] = useState("")

  const [websocket, setWebsocket] = useState(null)


  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8090/chat?roomId=${1}`);
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
  
  const sendMessage = (value) => {  
    console.log(value)
    websocket.send(value) 
    const newArr = [...messages.concat(value)] 
    setMessages(newArr)
  }
  
  if (websocket) {
    websocket.onmessage = (ev) => {
      console.log(ev.data)
    }
  }

  
  return (
  <>
    <div>
      {messages.map((el,i) => 
      (<p key={i}>{el}</p>)
      )}
    </div>
    <button onClick={() => sendMessage(inputValue)}>Send</button>
    <input type="text" onChange={(e) => (setInputValue(e.target.value))} />
  
  </>
  )
}

export default ChatRoom
import { useState } from "react"



function Chat({ws, chatName, msgHandler, messages}){
  console.log("chat changed to", chatName)
  const [inputValue, setInputValue] = useState("")
  console.log(messages)
  const sendMessage = (value) => {  
    const json_message = {
      "type": "NEW_MESSAGE",
      "body": `${value}`,
      "chatid": `${chatName}`
    }
    ws.send(JSON.stringify(json_message))
    msgHandler(chatName, value)
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

export default Chat
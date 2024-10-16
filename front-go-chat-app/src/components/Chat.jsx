import { useState } from "react"



function Chat({ws, chatName, msgHandler, messages}){
  console.log("chat changed to", chatName)

  const [inputValue, setInputValue] = useState("")
  console.log(messages)
  
  return (
  <>
 <div>
      {messages.map((el,i) => 
      (<p key={i}>{el}</p>)
      )}
    </div>
    <button onClick={() => msgHandler(chatName, inputValue)}>Send</button>
    <input type="text" onChange={(e) => (setInputValue(e.target.value))} />
  </>
  )
}

export default Chat
import { useState } from "react"



function Chat({ chat, msgHandler, subscribeHandler, messages}){
  console.log(messages)
  console.log("chat changed to", chat)
  const [inputValue, setInputValue] = useState("")
  console.log(messages)
  
  return (
  <>
    {
      chat.chat_type == 'group' && !chat.participation ? (
        <button key={chat.chat_id} onClick= {() => {subscribeHandler(chat.chat_id) }}> Join {chat.chat_id} group chat </button>
      ) : (
        <div>
        {messages.map((el,i) => 
          (<p key={i}>{el.username}: {el.body}</p>)
        )}
      <button onClick={() => msgHandler(chat, inputValue)}>Send</button>
      <input type="text" onChange={(e) => (setInputValue(e.target.value))} />
      </div>
    )
  }
    
     
      </>
    )
}

export default Chat
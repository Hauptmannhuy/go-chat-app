import { useContext, useState } from "react"
import { GlobalContext } from "../contexts/GlobalContext"



function Chat({ onSend}){
  const {selectedChat, messages} = useContext(GlobalContext)
  const chatMessages = messages[selectedChat.chat_name]
  console.log("chat changed to", selectedChat)
  const [inputValue, setInputValue] = useState("")
  
  return (
  <>
    {
      selectedChat.chat_type == 'group' && !selectedChat.participation ? (
        <button key={selectedChat.chat_id} onClick= {() => {subscribeHandler(selectedChat.chat_id) }}> Join {selectedChat.chat_id} group chat </button>
      ) : (
        chatMessages ? ( <div className="messages-display">
          {chatMessages.map((el,i) => 
            (<p key={i}>{el.username}: {el.body}</p>)
          )}
        </div>) : (null)
       
    )
  }
    
  <button onClick={() => onSend(selectedChat, inputValue)}>Send</button>
  <input type="text" onChange={(e) => (setInputValue(e.target.value))} />
     
      </>
    )
}

export default Chat
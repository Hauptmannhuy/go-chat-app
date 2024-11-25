import { ChatSnippet } from "./ChatSnippet"



function ChatList({chats, handleSelect, handleJoin, currentUsername, messages}) {
  const chatKeys = Object.keys(chats)

  return (
  <>

   <h2>Chat list</h2>
    <div className="chat-list">

    {chatKeys.map((key) => (
        <div className="chat-snippet" onClick={ () => {handleSelect(key) }}>
          { 
            < ChatSnippet 
              name={key}
              messages={messages}
              currentUsername={currentUsername}/> 
          }</div>
      )

        
    )}
      
    </div>
  
  </>
  )
}



export default ChatList
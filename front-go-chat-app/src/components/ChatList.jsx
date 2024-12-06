import { ChatSnippet } from "./ChatSnippet"



function ChatList({onSelect}) {
  
  const chatKeys = Object.keys(chats)

  return (
  <>
   <h2>Chat list</h2>
    <div className="chat-list">

    {chatKeys.map((key) => (
        <div className="chat-snippet" onClick={ () => {onSelect(key) }}>
          { 
            < ChatSnippet 
              name={key}
              /> 
          }</div>
      )

        
    )}
    </div>

   
  
  </>
  )
}



export default ChatList
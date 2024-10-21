
function ChatList({chats, handleSelect, handleJoin}) {
  console.log(chats)

  return (
  <>

   <h2>Chat list</h2>
    <div className="chat-list">

    {chats.map((chat) => (
       chat.participation ? ( 
        <button  

          key={chat.name} 
          onClick= {() => {handleSelect(chat.name) }}
          name={chat.name}
          >{chat.name} 

        </button>)
         : 

        (
        <button
        key={chat.name}
        onClick = {() => {handleJoin(chat.name)}}>
          Join chat {chat.name}
        </button>
        )

       
      ))}
      
    </div>
  
  </>
  )
}



export default ChatList
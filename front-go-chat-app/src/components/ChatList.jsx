
function ChatList({chats, handleSelect, handleJoin}) {
  

  return (
  <>

   <h2>Chat list</h2>
    <div className="chat-list">

    {chats.map((el) => (
       el.participation ? ( 
        <button  

          key={el.name} 
          onClick= {() => {handleSelect(el.name) }}
          name={el.name}
          >{el.name} 

        </button>)
         : 

        (
        <button
        key={el.name}
        onClick = {() => {handleJoin(el.name)}}>
          Join chat {el.name}
        </button>
        )

       
      ))}
      
    </div>
  
  </>
  )
}



export default ChatList
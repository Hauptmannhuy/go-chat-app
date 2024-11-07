
function ChatList({chats, handleSelect, handleJoin}) {
  const chatKeys = Object.keys(chats)
  console.log(chats)
  return (
  <>

   <h2>Chat list</h2>
    <div className="chat-list">

    {chatKeys.map((key) => (
       
        <button key={key} onClick= {() => {handleSelect(key) }} name={key}> {chats[key].name } </button>)
    )}
      
    </div>
  
  </>
  )
}



export default ChatList
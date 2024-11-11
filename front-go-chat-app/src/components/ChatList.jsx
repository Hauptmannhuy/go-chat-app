
function ChatList({chats, handleSelect, handleJoin, currentUsername}) {
  const chatKeys = Object.keys(chats)
  console.log(chats)
  return (
  <>

   <h2>Chat list</h2>
    <div className="chat-list">

    {chatKeys.map((key) => (
       
        <button key={key} onClick= {() => {handleSelect(key) }} name={key}> { key.split('_').filter((el) => (el != currentUsername)).join(' ') } </button>)
    )}
      
    </div>
  
  </>
  )
}



export default ChatList
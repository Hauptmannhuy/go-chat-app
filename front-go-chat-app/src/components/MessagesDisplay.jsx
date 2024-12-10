export function MessagesDisplay({chatMessages}){
 
  return  (
    <>
    {
      chatMessages ? ( <div className="messages-display">
        {chatMessages.map((el,i) => 
          (<p key={i}>{el.username}: {el.body}</p>)
        )}
      </div>) : (null)
    }
    </>
  )
}
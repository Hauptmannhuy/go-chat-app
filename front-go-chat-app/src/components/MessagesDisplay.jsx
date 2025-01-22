import Message from "./Message"
import "../assets/message.css"

export function MessagesDisplay({chatMessages}){
 console.log(chatMessages)
  return  (
    <>
    {
      chatMessages ? ( <div className="messages-display">
        {chatMessages.map((el,i) => (
          <Message
          key={i}
          body={el.body}
          username={el.username}/>)
        )}
      </div>) : (null)
    }
    </>
  )
}
import { useAuth } from "../modules/useAuth"

export default function Message({username, body}) {
  const {getUsername} = useAuth()
  
  if (getUsername() == username) {
    return (
      <div className="message">
        <p className="message-text">{username}: {body}</p>
      </div>
    )
  } else {
    return (
      <div className="message not-user">
      <p className="message-text">{username}: {body}</p>
    </div>
    )
  }
}
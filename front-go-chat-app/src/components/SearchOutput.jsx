export default function SearchOutput({outputData}){
  const chats = outputData

  return (<>
  
    {chats.map((chat) => (
      <p>{chat.name}</p>
    ))}
  
  </>)
}
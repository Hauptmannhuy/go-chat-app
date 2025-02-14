export const useEnvelope = () => {
  const makeEnvelope = (type, data) => {
    const actions = {
      "NEW_MESSAGE": {
        type: "NEW_MESSAGE",
        chat_name: `${data[0]}`,
        body: `${data[1]}`,
      },
      "NEW_CHAT": {
        type: "NEW_CHAT",
        chat_name: `${data[0]}`,
      },
      "NEW_PRIVATE_CHAT": {
        type: "NEW_PRIVATE_CHAT",
        receiver_id: Number(`${data[0]}`),
        body: `${data[1]}`,
      },
      "NEW_GROUP_CHAT":{
        type:"NEW_GROUP_CHAT",
        chat_name: `${data[0]}`,
      },
      "JOIN_CHAT": {
        type: "JOIN_CHAT",
        chat_id: Number(`${data[0]}`),
        chat_name: `${data[1]}`,
        body_message:`${data[2]}`,
      },
      "SEARCH_QUERY": {
        type: "SEARCH_QUERY",
        input: `${data[0]}`,
      }
    }

    const json = actions[type];
    if (json) {
      return json
    } else {
      console.error(`Unknown action type: ${type}`);
    }
  };
  return {makeEnvelope}
}
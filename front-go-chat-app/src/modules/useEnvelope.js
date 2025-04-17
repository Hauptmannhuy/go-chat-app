import { FileImage } from "../types/image";

export const useEnvelope = () => {


  
  /**
   * 
   * @param {String} type
   * @param {{input: string, image: FileImage}} payload 
   * @param {Array<any>} info
   */
  const makeEnvelope = (type, payload, info) => {
    console.log(payload.image)
    const actions = {
      "NEW_MESSAGE": {
        type: "NEW_MESSAGE",
        chat_name: `${info[0]}`,
        body: `${payload.input}`,
        image: payload.image
      },
      "NEW_CHAT": {
        type: "NEW_CHAT",
        chat_name: `${info[0]}`,
      },
      "NEW_PRIVATE_CHAT": {
        type: "NEW_PRIVATE_CHAT",
        receiver_id: Number(`${info[0]}`),
        body: `${payload.input}`,
        image: payload.image
      },
      "NEW_GROUP_CHAT":{
        type:"NEW_GROUP_CHAT",
        chat_name: `${info[0]}`,
      },
      "JOIN_CHAT": {
        type: "JOIN_CHAT",
        chat_id: Number(`${info[0]}`),
        chat_name: `${info[1]}`,
        body_message: payload.input,
        body_image: payload.image
      },
      "SEARCH_QUERY": {
        type: "SEARCH_QUERY",
        input: payload.input,
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
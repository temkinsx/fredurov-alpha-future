import type { Chat } from "../../../types.ts";
import {useState} from "react";

const mockChats: Chat[] = [
    {
        name: "Финансовая аналитика",
        date: "10:30"
    },
    {
        name: "Проверка платежей",
        date: "Вчера"
    },
    {
        name: "Маркетинговая кампания",
        date: "2 дня назад"
    }
]

function ChatsList() {
    const [activeChat, setActiveChat] = useState<string>(mockChats[0].name);

    const handleChangeChat = (chatName: string) => {
        setActiveChat(chatName);
    }

    return (
        <div className="p-3 overflow-y-scroll">
            <h2 className="text-gray-500 text-sm uppercase">Недавние чаты</h2>
            {mockChats.map(chat => (
                <div key={ chat.name }
                    className={`${activeChat === chat.name ? "bg-red-100 outline-1 outline-red-200" : ""} mt-5 p-2 hover:bg-red-50 
                    rounded-2xl transition-all outline-red-200`}
                        onClick={() => handleChangeChat(chat.name) }
                >
                    <h3>{ chat.name }</h3>
                    <p className="text-gray-400 text-sm">{ chat.date }</p>
                </div>
            ))}
        </div>
    );
}

export default ChatsList;
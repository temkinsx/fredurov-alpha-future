import { Plus } from "lucide-react";
import ChatsList from "./ChatsList.tsx";


function ChatAside() {

    const handleCreateNewChat = () => {

    }

    return (
        <aside className="flex flex-col border-r border-r-gray-100" style={{gridColumn: "1/2", gridRow: "1/11"}}>
            <div className="flex justify-center items-center border-b border-b-gray-100 h-30">
                <button
                    className="flex justify-center items-center
                gap-2 text-white bg-red-500 p-3.5 w-8/10
                rounded-2xl h-15 hover:bg-red-600 transition-all
                shadow-xl shadow-red-200"
                    onClick={ handleCreateNewChat }
                >
                    <Plus className="h-6 w-6"/>
                    Новый чат
                </button>
            </div>
            <ChatsList />
        </aside>
    );
}

export default ChatAside;
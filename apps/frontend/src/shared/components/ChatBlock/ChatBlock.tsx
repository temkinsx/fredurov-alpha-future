import QuickActions from "./QuickActions/QuickActions.tsx";
import MainChat from "./MainChat/MainChat.tsx";

function ChatBlock() {

    return (
        <main
            className="bg-gray-50 p-6 overflow-y-scroll"
            style={{gridColumn: "2/3", gridRow: "2/10"}}
        >
            <QuickActions />
            <MainChat />
        </main>
    );
}

export default ChatBlock;
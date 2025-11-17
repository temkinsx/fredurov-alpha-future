import ChatHeader from "../../shared/components/ChatHeader/ChatHeader.tsx";
import ChatAside from "../../shared/components/ChatAside/ChatAside.tsx";
import styles from "./ChatPage.module.css";
import ChatBlock from "../../shared/components/ChatBlock/ChatBlock.tsx";
import QueryInput from "../../shared/components/ChatBlock/QueryInput/QueryInput.tsx";


function ChatPage() {

    return (
        <div className={ styles.pageContainer }>
            <ChatAside />
            <ChatHeader />
            <ChatBlock />
            <QueryInput />
        </div>
    );
}

export default ChatPage;
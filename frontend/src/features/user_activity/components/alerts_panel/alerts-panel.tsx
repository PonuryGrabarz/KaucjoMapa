import Styles from "./alerts-panel.module.sass";
import AlertsList from "../alerts_list/alerts-list";
import AlertsSummary from "../alerts_summary/alerts-summary";

import { useState } from "preact/hooks";

import type { AlertItem } from "../../lib/alerts-mock-data";

interface Props {
    items: AlertItem[];
}

const ALERTS_LIST_ID = "alerts-list";

export default function AlertsPanel({ items }: Props) {
    const [allRead, setAllRead] = useState(false);
    const initialUnreadCount = items.filter(({ unread }) => unread).length;
    const unreadCount = allRead ? 0 : initialUnreadCount;

    return (
        <div class={Styles["alerts-panel"]} data-all-read={allRead ? "true" : "false"}>
            <AlertsSummary
                listId={ALERTS_LIST_ID}
                unreadCount={unreadCount}
                onMarkAllRead={() => setAllRead(true)}
            />

            <AlertsList id={ALERTS_LIST_ID} items={items} allRead={allRead} />
        </div>
    );
}

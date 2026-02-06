import { Contact, Inbox, LucideProps, Settings } from "lucide-react"

export interface NavMain {
    title: string
    isActive: boolean
    icon: React.ForwardRefExoticComponent<Omit<LucideProps, "ref"> & React.RefAttributes<SVGSVGElement>>
}

export const NAV_TITLE_CHAT = "Chat" , 
NAV_TITLE_FRIEND = "Friend" ,
NAV_TITLE_Setting = "Setting" 

export const Navigation: NavMain[] = [
    {
        title: "Chat",
        icon: Inbox,
        isActive: true
    },
    {
        title: "Friend",
        icon: Contact,
        isActive: true
    },
    {
        title: "Setting",
        icon: Settings,
        isActive: true
    }
]
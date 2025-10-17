import { Icon } from "@iconify/react";
import { toast } from "sonner";

export function ToastSuccess(message: string) {
  Toast(message , "ep:success-filled", "#0BE670")
}

export function ToastError(message: string) {
  Toast(message , "material-symbols:error", "#e12626")
}

function Toast(message: string, icon: string, colorIcon: string) {
  toast(message , {
    icon: <Icon icon={icon} color={colorIcon}/>,
    position: 'top-center'
  })
}


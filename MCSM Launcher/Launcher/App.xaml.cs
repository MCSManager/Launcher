using System;
using System.Collections.Generic;
using System.Configuration;
using System.Data;
using System.Linq;
using System.Threading.Tasks;
using System.Windows;

namespace Launcher
{
    /// <summary>
    /// App.xaml 的交互逻辑
    /// </summary>
    public partial class App : Application
    {
        private bool IsSupportDarkMode = Environment.OSVersion.Version.Major == 10 &&
            Environment.OSVersion.Version.Build >= 17763;

        protected override void OnStartup(StartupEventArgs e)
        {
            // 深色模式 兼容
            if (IsSupportDarkMode)
            {
                UxTheme.AllowDarkModeForApp(true);
                UxTheme.ShouldSystemUseDarkMode();
            }
        }
    }
}

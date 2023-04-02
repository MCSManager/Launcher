using ModernWpf.Controls;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Text.Json.Nodes;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;



namespace Launcher
{
    /// <summary>
    /// MainWindow.xaml 的交互逻辑
    /// </summary>
    public partial class MainWindow : Window
    {
        public MainWindow()
        {
            InitializeComponent();
        }

        private static Hardcodet.Wpf.TaskbarNotification.TaskbarIcon TaskbarIcon { get; set; } = new()
        {
            IconSource = new BitmapImage(new Uri("pack://application:,,,/Launcher;component/icon.ico")),
            ToolTipText = "MCSM已挂至后台，单击显示",
        };

        protected override async void OnInitialized(EventArgs e)
        {
            string file = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager","web","data", "SystemConfig","config.json");
            if (File.Exists(file))
            {
                using var reader = new StreamReader(file);
                var node = JsonNode.Parse(await reader.ReadToEndAsync());

                MS_Run_Port.Text = node["httpPort"].ToString();
                WWwww.Tag = $"http://localhost:{node["httpPort"]}";
            }
            else
            {
#if !DEBUG
                MessageBox.Show("找不到配置文件，试试以管理员权限打开启动器  / 重新下载。", "LauncehrWrapper", MessageBoxButton.OK, MessageBoxImage.Hand);
                Environment.Exit(0);
                return;
#endif
            }

            MS_Run_State.Text = "关闭";

            var menui1 = new MenuItem() { Header = "退出" };
            menui1.Click += (sender, args) =>
            {
                processes.ForEach(process =>
                {
                    if (!process.HasExited) process.Kill();
                });
                App.Current.Shutdown();
            };

            TaskbarIcon.ContextMenu = new ContextMenu()
            {
                Items =
                {
                    menui1
                }
            };
            TaskbarIcon.TrayLeftMouseUp += (sender, args) =>
            {
                Activate();
                Visibility = Visibility.Visible;
            };

            base.OnInitialized(e);
        }
        private static List<Process> processes { get; set; } = new();

        protected override void OnClosing(CancelEventArgs e)
        {
            e.Cancel = true;
            Visibility = Visibility.Hidden;
        }





        private void HyperlinkButton_Click_1(object sender, RoutedEventArgs e)
        {
            try
            {
                Process.Start(Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "daemon", "logs", "current.log"));
            }
            catch { }
            try
            {
                Process.Start(Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "web", "logs", "current.log"));
            }
            catch { }
        }



        /// <summary>
        /// 状态
        /// </summary>
        private static bool State { get; set; }

        /// <summary>
        /// 启动
        /// </summary>
        /// <param name="sender"></param>
        /// <param name="e"></param>
        private void Button_Click(object sender, RoutedEventArgs e)
        {
            var button = (sender as Button);
            if (!State)
            {
                string daemonPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "daemon");
                string webPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "mcsmanager", "web");
                var p1 = Process.Start(new ProcessStartInfo(Path.Combine(daemonPath, "node_app.exe"), $"\"{Path.Combine(daemonPath, "app.js")}\"")
                {
                    CreateNoWindow = true,
                    UseShellExecute = false,
                    WorkingDirectory = daemonPath
                });
                var p2 = Process.Start(new ProcessStartInfo(Path.Combine(webPath, "node_app.exe"), $"\"{Path.Combine(webPath, "app.js")}\"")
                {
                    CreateNoWindow = true,
                    UseShellExecute = false,
                    WorkingDirectory = webPath,
                });
                p1.EnableRaisingEvents = p2.EnableRaisingEvents = true;
                if (!p1.HasExited && !p2.HasExited)
                {
                    processes.Add(p1);
                    processes.Add(p2);
                    p1.Exited += (s, e) => 
                    {
                        State = false;
                        App.Current.Dispatcher.Invoke(() =>
                        {
                            button.Content = "开启后台程序";
                        });
                        if (!p1.HasExited) p1.Kill();
                        if (!p2.HasExited) p2.Kill();
                        processes.Clear();
                    };
                    p2.Exited += (s, e) =>
                    {
                        State = false;
                        App.Current.Dispatcher.Invoke(() =>
                        {
                            button.Content = "开启后台程序";
                        });
                        if (!p1.HasExited) p1.Kill();
                        if (!p2.HasExited) p2.Kill();
                        processes.Clear();
                    };
                    button.Content = "关闭后台程序";
                    State = true;
                    MS_Run_State.Text = "开启";
                    return;
                }
                button.Content = "关闭后台程序";
                State = true;
                if (!p1.HasExited) p1.Kill();
                if (!p2.HasExited) p2.Kill();
            }
            else
            {
                processes.ToList().ForEach(process =>
                {
                    process.EnableRaisingEvents = false;
                    if (!process.HasExited) process.Kill();
                });
                MS_Run_State.Text = "关闭";
                State = false;
                button.Content = "开启后台程序";
            }

        }


        /// <summary>
        /// 打开面板连接
        /// </summary>
        private void Button_Click_1(object sender, RoutedEventArgs e)
        {
            try
            {
                Process.Start((sender as Button).Tag.ToString());
            }
            catch
            {
                MessageBox.Show("貌似没有读到配置呢，试试重新打开启动器。", "LauncehrWrapper", MessageBoxButton.OK, MessageBoxImage.Hand);
            }
            
        }
    }
}

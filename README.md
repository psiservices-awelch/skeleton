DESCRIPTION:
    Skeleton copies an entire directory and all it's sub directories to a new location. It is meant for copying a base project with like components to a new location to make the setup of a new project with like elements easier to start.

    Skeleton checks that the directory that you wish to copy from exists and will error if it does not.

    Skeleton checks that the desired directory you wish to copy to does not exist before starting it's work, if it does exist it warns the user that the directory exists and that the current contents will be removed before asking if they wish to continue. If the user wishes to continue they can enter yes or y, capitalization does not matter, any other input aborts the program. If the user continues the current directory will be removed and recreated. If the directory does not exist it will be created.

    After Skeleton creates the directory to copy to it initalizes git to turn it into a git repo before begining the copy process. ***This means that git must be installed for you to use this tool!!!***

    Skeletor then generates a list of all the files and directories in the source project and begins the copy process. If the copied files are go source files, go.mod files, or go.sum files then the program also offers, if the -r option was entered when running the command, the ability to replace the old project name name with the name of the project you are going to create. 
            For Example the file:

                Old file:
                    project main

                    import github.com/awelch/microservice_base/server

            could become:

                New file:
                    project main

                    import github.com/awelch/test_service/server

    There is one small danger that should be noted, Skeleton goes through the whole file when it checks for the old name to replace so the name of the old project should be as unique as possible if you use this feature or variables, functions, comments, etc could be affected. You may want to manually update the copied files if you use packages or other code that have a similar name to the project you are copying.

USING SKELETON:
    In order to create the binary to use Skeleton you can simply download the source code and run 'go build' in the CLI. The 'skeleton' binary should appear, for linux, or the skeleton.exe file, for windows, and be ready for use.

    The base Usage of Skeleton is the following:

    Usage: ./skeleton [options] <SOURCE PATH> <DEST PATH>:
    Options:
      -r string
        the old and new name to be used to update go files imports separated by a colon
      -h
        views the skeleton command help

    The SOURCE PATH and DEST PATH are required in order for the program to run.
